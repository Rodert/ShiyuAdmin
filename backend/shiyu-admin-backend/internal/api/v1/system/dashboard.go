package system

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

type geoAccessPointVO struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Count       int     `json:"count"`
}

type geoAccessTopIPVO struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
	Count   int    `json:"count"`
}

type geoAccessSummaryVO struct {
	WindowHours int                 `json:"window_hours"`
	Total       int                 `json:"total"`
	Countries   []*geoAccessPointVO `json:"countries"`
	TopIPs      []*geoAccessTopIPVO `json:"top_ips"`
	UpdatedAt   int64               `json:"updated_at"`
}

type geoMeta struct {
	Country     string
	CountryCode string
	City        string
	Lat         float64
	Lon         float64
}

type geoCacheEntry struct {
	Meta      geoMeta
	ExpiresAt time.Time
}

type ipAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

var (
	geoLookupHTTPClient = &http.Client{Timeout: 2 * time.Second}
	geoCache            sync.Map
)

func registerDashboardRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, operationLogSvc interfaces.OperationLogService) {
	if operationLogSvc == nil {
		return
	}
	rg.GET("/dashboard/geo-access", middleware.RequirePermission(permissionSvc, "dashboard:view"), func(c *gin.Context) {
		getGeoAccessSummary(c, operationLogSvc)
	})
}

func getGeoAccessSummary(c *gin.Context, operationLogSvc interfaces.OperationLogService) {
	window := 24
	if rawWindow := strings.TrimSpace(c.Query("hours")); rawWindow != "" {
		var parsed int
		_, _ = fmt.Sscanf(rawWindow, "%d", &parsed)
		if parsed > 0 && parsed <= 24*30 {
			window = parsed
		}
	}

	start := time.Now().Add(-time.Duration(window) * time.Hour).Unix()
	pageResult, err := operationLogSvc.Page(c, &dto.OperationLogPageRequest{
		Page:      1,
		PageSize:  100,
		StartTime: &start,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if pageResult == nil || len(pageResult.Items) == 0 {
		response.Success(c, &geoAccessSummaryVO{
			WindowHours: window,
			Total:       0,
			Countries:   []*geoAccessPointVO{},
			TopIPs:      []*geoAccessTopIPVO{},
			UpdatedAt:   time.Now().Unix(),
		})
		return
	}

	type countryBucket struct {
		Name  string
		Code  string
		Lat   float64
		Lon   float64
		Count int
	}
	type ipBucket struct {
		IP      string
		Country string
		City    string
		Count   int
	}

	countryMap := make(map[string]*countryBucket)
	ipMap := make(map[string]*ipBucket)
	total := 0

	for _, item := range pageResult.Items {
		if item == nil {
			continue
		}
		ip := strings.TrimSpace(item.IP)
		if ip == "" || ip == "unknown" {
			continue
		}
		meta := resolveGeoByIP(ip)
		total++

		countryKey := meta.CountryCode
		if countryKey == "" {
			countryKey = "ZZ"
		}
		cb, exists := countryMap[countryKey]
		if !exists {
			cb = &countryBucket{
				Name: meta.Country,
				Code: meta.CountryCode,
				Lat:  meta.Lat,
				Lon:  meta.Lon,
			}
			countryMap[countryKey] = cb
		}
		cb.Count++

		ib, exists := ipMap[ip]
		if !exists {
			ib = &ipBucket{
				IP:      ip,
				Country: meta.Country,
				City:    meta.City,
			}
			ipMap[ip] = ib
		}
		ib.Count++
	}

	countryItems := make([]*geoAccessPointVO, 0, len(countryMap))
	for _, x := range countryMap {
		countryItems = append(countryItems, &geoAccessPointVO{
			Country:     x.Name,
			CountryCode: x.Code,
			City:        "",
			Lat:         x.Lat,
			Lon:         x.Lon,
			Count:       x.Count,
		})
	}
	sortGeoPoints(countryItems)

	ipItems := make([]*geoAccessTopIPVO, 0, len(ipMap))
	for _, x := range ipMap {
		ipItems = append(ipItems, &geoAccessTopIPVO{
			IP:      x.IP,
			Country: x.Country,
			City:    x.City,
			Count:   x.Count,
		})
	}
	sortTopIPs(ipItems)
	if len(ipItems) > 8 {
		ipItems = ipItems[:8]
	}

	response.Success(c, &geoAccessSummaryVO{
		WindowHours: window,
		Total:       total,
		Countries:   countryItems,
		TopIPs:      ipItems,
		UpdatedAt:   time.Now().Unix(),
	})
}

func sortGeoPoints(items []*geoAccessPointVO) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Count > items[i].Count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func sortTopIPs(items []*geoAccessTopIPVO) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Count > items[i].Count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func resolveGeoByIP(ip string) geoMeta {
	cached := readGeoCache(ip)
	if cached != nil {
		return *cached
	}

	meta := buildFallbackGeo(ip)
	if parsed, ok := parseIP(ip); ok && !isPrivateOrLoopback(parsed) {
		if remote, err := queryIPAPI(ip); err == nil {
			meta = remote
		}
	}
	writeGeoCache(ip, meta)
	return meta
}

func parseIP(ip string) (netip.Addr, bool) {
	p, err := netip.ParseAddr(ip)
	if err == nil {
		return p, true
	}
	if host, _, splitErr := net.SplitHostPort(ip); splitErr == nil {
		p2, err2 := netip.ParseAddr(host)
		return p2, err2 == nil
	}
	return netip.Addr{}, false
}

func isPrivateOrLoopback(addr netip.Addr) bool {
	return addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() || addr.IsMulticast()
}

func buildFallbackGeo(ip string) geoMeta {
	if parsed, ok := parseIP(ip); ok && isPrivateOrLoopback(parsed) {
		return geoMeta{
			Country:     "Local Network",
			CountryCode: "LOCAL",
			City:        "Private",
			Lat:         39.9042,
			Lon:         116.4074,
		}
	}
	return geoMeta{
		Country:     "Unknown",
		CountryCode: "ZZ",
		City:        "Unknown",
		Lat:         0,
		Lon:         0,
	}
}

func queryIPAPI(ip string) (geoMeta, error) {
	endpoint := "http://ip-api.com/json/" + url.PathEscape(ip) + "?fields=status,country,countryCode,city,lat,lon"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return geoMeta{}, err
	}
	resp, err := geoLookupHTTPClient.Do(req)
	if err != nil {
		return geoMeta{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return geoMeta{}, fmt.Errorf("ip-api status %d", resp.StatusCode)
	}
	var body ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return geoMeta{}, err
	}
	if body.Status != "success" {
		return geoMeta{}, fmt.Errorf("ip-api failure")
	}
	return geoMeta{
		Country:     strings.TrimSpace(body.Country),
		CountryCode: strings.ToUpper(strings.TrimSpace(body.CountryCode)),
		City:        strings.TrimSpace(body.City),
		Lat:         body.Lat,
		Lon:         body.Lon,
	}, nil
}

func readGeoCache(ip string) *geoMeta {
	v, ok := geoCache.Load(ip)
	if !ok {
		return nil
	}
	entry, ok := v.(geoCacheEntry)
	if !ok {
		return nil
	}
	if time.Now().After(entry.ExpiresAt) {
		geoCache.Delete(ip)
		return nil
	}
	meta := entry.Meta
	return &meta
}

func writeGeoCache(ip string, meta geoMeta) {
	geoCache.Store(ip, geoCacheEntry{
		Meta:      meta,
		ExpiresAt: time.Now().Add(6 * time.Hour),
	})
}
