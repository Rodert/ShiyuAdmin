import { request } from '@umijs/max';

export interface GeoCountryPoint {
  country: string;
  country_code: string;
  city: string;
  lat: number;
  lon: number;
  count: number;
}

export interface GeoTopIP {
  ip: string;
  country: string;
  city: string;
  count: number;
}

export interface GeoAccessSummary {
  window_hours: number;
  total: number;
  countries: GeoCountryPoint[];
  top_ips: GeoTopIP[];
  updated_at: number;
}

export async function getGeoAccessSummary(
  params?: { hours?: number },
  opts?: { skipErrorHandler?: boolean },
) {
  return request<{
    code: number;
    data: GeoAccessSummary;
    message?: string;
  }>('/api/v1/system/dashboard/geo-access', {
    method: 'GET',
    params,
    ...opts,
  });
}
