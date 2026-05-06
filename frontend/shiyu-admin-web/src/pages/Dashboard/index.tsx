import { FullscreenExitOutlined, FullscreenOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import type { MenuDataItem } from '@ant-design/pro-components';
import { request, useModel } from '@umijs/max';
import { Button, Col, Row, Space, Spin, Table, Tag, Typography } from 'antd';
import type { EChartsOption } from 'echarts';
import * as echarts from 'echarts';
import dayjs from 'dayjs';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { GeoAccessSummary } from '@/services/shiyu-api/dashboard';
import { getGeoAccessSummary } from '@/services/shiyu-api/dashboard';
import type { CacheStats, DatabaseStats, OnlineUser } from '@/services/shiyu-api/monitor';
import { getCacheStats, getDatabaseStats, getOnlineUsers } from '@/services/shiyu-api/monitor';
import type { OperationLog } from '@/services/shiyu-api/operation_log';
import { getOperationLogList } from '@/services/shiyu-api/operation_log';
import { getDeptList } from '@/services/shiyu-api/dept';
import { getRoleList } from '@/services/shiyu-api/role';
import { getUserList } from '@/services/shiyu-api/user';
import './dashboard.less';

const countMenuItems = (menus?: MenuDataItem[]): number => {
  if (!menus?.length) {
    return 0;
  }
  return menus.reduce((total, item) => total + 1 + countMenuItems(item.children), 0);
};

interface SystemHealthResponse {
  code: number;
  data?: { status?: string; time?: number };
}

const TIME_SLOTS = 12;

function bucketLogsByHour(logs: OperationLog[], slotCount: number): number[] {
  const buckets = new Array(slotCount).fill(0);
  const now = Date.now() / 1000;
  const slotSec = (24 * 3600) / slotCount;
  for (const log of logs) {
    const age = now - log.created_at;
    if (age < 0 || age > 24 * 3600) {
      continue;
    }
    const idx = Math.min(slotCount - 1, Math.floor(age / slotSec));
    buckets[slotCount - 1 - idx]++;
  }
  return buckets;
}

/** 避免 echarts-for-react 与 React 19 类型不兼容，使用原生 init */
function ChartShell({
  option,
  style,
}: {
  option: EChartsOption;
  style?: React.CSSProperties;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const chartRef = useRef<ReturnType<typeof echarts.init> | null>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) {
      return;
    }
    const chart = echarts.init(el);
    chartRef.current = chart;
    const ro = new ResizeObserver(() => chart.resize());
    ro.observe(el);
    return () => {
      ro.disconnect();
      chart.dispose();
      chartRef.current = null;
    };
  }, []);

  useEffect(() => {
    chartRef.current?.setOption(option, true);
  }, [option]);

  return <div ref={ref} style={style} />;
}

const Dashboard: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const userName = currentUser?.name || '管理员';

  const [loading, setLoading] = useState(true);
  const [healthOk, setHealthOk] = useState<boolean | null>(null);
  const [userTotal, setUserTotal] = useState<number | null>(null);
  const [roleTotal, setRoleTotal] = useState<number | null>(null);
  const [deptCount, setDeptCount] = useState<number | null>(null);
  const [logTotal, setLogTotal] = useState<number | null>(null);
  const [logsSample, setLogsSample] = useState<OperationLog[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<OnlineUser[]>([]);
  const [cacheStats, setCacheStats] = useState<CacheStats | null>(null);
  const [dbStats, setDbStats] = useState<DatabaseStats | null>(null);
  const [geoSummary, setGeoSummary] = useState<GeoAccessSummary | null>(null);
  const [worldMapReady, setWorldMapReady] = useState(false);

  const boardRef = useRef<HTMLDivElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);

  useEffect(() => {
    const sync = () => {
      setIsFullscreen(document.fullscreenElement === boardRef.current);
    };
    document.addEventListener('fullscreenchange', sync);
    return () => document.removeEventListener('fullscreenchange', sync);
  }, []);

  const toggleFullscreen = useCallback(async () => {
    const el = boardRef.current;
    if (!el) {
      return;
    }
    try {
      if (!document.fullscreenElement) {
        await el.requestFullscreen();
      } else {
        await document.exitFullscreen();
      }
    } catch (_e) {
      // 部分浏览器需用户手势或策略限制，静默失败
    }
  }, []);

  const loadData = useCallback(async () => {
    setLoading(true);
    const skip = { skipErrorHandler: true } as const;

    const healthP = request<SystemHealthResponse>('/api/v1/system/health', { method: 'GET', ...skip }).then(
      (r) => (r.code === 200 && r.data?.status === 'ok' ? true : false),
      () => null,
    );

    const usersP = getUserList({ page: 1, page_size: 1 }, skip).then(
      (r) => (r.code === 200 ? r.data?.total ?? null : null),
      () => null,
    );

    const rolesP = getRoleList({ page: 1, page_size: 1 }, skip).then(
      (r) => (r.code === 200 ? r.data?.total ?? null : null),
      () => null,
    );

    const deptsP = getDeptList(skip).then(
      (r) => (r.code === 200 && Array.isArray(r.data) ? r.data.length : null),
      () => null,
    );

    const logsMetaP = getOperationLogList({ page: 1, page_size: 1 }, skip).then(
      (r) => (r.code === 200 ? r.data?.total ?? null : null),
      () => null,
    );

    const logsSampleP = getOperationLogList({ page: 1, page_size: 120 }, skip).then(
      (r) => (r.code === 200 ? r.data?.items ?? [] : []),
      () => [],
    );

    const onlineP = getOnlineUsers(skip).then(
      (r) => (r.code === 200 ? r.data ?? [] : []),
      () => [],
    );

    const cacheP = getCacheStats(skip).then(
      (r) => (r.code === 200 ? r.data : null),
      () => null,
    );

    const dbP = getDatabaseStats(skip).then(
      (r) => (r.code === 200 ? r.data : null),
      () => null,
    );
    const geoP = getGeoAccessSummary({ hours: 24 }, skip).then(
      (r) => (r.code === 200 ? r.data : null),
      () => null,
    );

    const results = await Promise.all([
      healthP,
      usersP,
      rolesP,
      deptsP,
      logsMetaP,
      logsSampleP,
      onlineP,
      cacheP,
      dbP,
      geoP,
    ]);

    setHealthOk(results[0] as boolean | null);
    setUserTotal(results[1] as number | null);
    setRoleTotal(results[2] as number | null);
    setDeptCount(results[3] as number | null);
    setLogTotal(results[4] as number | null);
    setLogsSample(results[5] as OperationLog[]);
    setOnlineUsers(results[6] as OnlineUser[]);
    setCacheStats(results[7] as CacheStats | null);
    setDbStats(results[8] as DatabaseStats | null);
    setGeoSummary(results[9] as GeoAccessSummary | null);
    setLoading(false);
  }, []);

  useEffect(() => {
    let mounted = true;
    fetch('https://fastly.jsdelivr.net/npm/echarts@5/map/json/world.json')
      .then((res) => res.json())
      .then((geojson) => {
        if (!mounted) {
          return;
        }
        echarts.registerMap('world', geojson as never);
        setWorldMapReady(true);
      })
      .catch(() => {
        if (mounted) {
          setWorldMapReady(false);
        }
      });
    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    loadData();
    const t = window.setInterval(loadData, 60000);
    return () => window.clearInterval(t);
  }, [loadData]);

  const menuCount = countMenuItems(initialState?.menuData);

  const hourly = useMemo(() => bucketLogsByHour(logsSample, TIME_SLOTS), [logsSample]);

  const timeAxisLabels = useMemo(
    () =>
      Array.from({ length: TIME_SLOTS }, (_, i) => {
        const hoursAgo = Math.round(((TIME_SLOTS - 1 - i) / TIME_SLOTS) * 24);
        if (hoursAgo <= 0) {
          return '现在';
        }
        return `${hoursAgo}h`;
      }),
    [],
  );

  const lineOption: EChartsOption = useMemo(
    () => ({
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(15,23,42,0.92)',
        borderColor: 'rgba(56,189,248,0.35)',
        textStyle: { color: '#e2e8f0' },
      },
      grid: { left: 44, right: 12, top: 32, bottom: 28 },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: timeAxisLabels,
        axisLine: { lineStyle: { color: 'rgba(148,163,184,0.3)' } },
        axisLabel: { color: 'rgba(148,163,184,0.75)', fontSize: 10, interval: 1, rotate: 28 },
      },
      yAxis: {
        type: 'value',
        splitLine: { lineStyle: { color: 'rgba(148,163,184,0.12)' } },
        axisLabel: { color: 'rgba(148,163,184,0.65)', fontSize: 10 },
      },
      series: [
        {
          name: '操作次数',
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 6,
          lineStyle: {
            width: 2,
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#22d3ee' },
              { offset: 1, color: '#818cf8' },
            ]),
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(56,189,248,0.35)' },
              { offset: 1, color: 'rgba(56,189,248,0)' },
            ]),
          },
          data: hourly,
        },
      ],
    }),
    [hourly, timeAxisLabels],
  );

  const geoMapOption: EChartsOption = useMemo(() => {
    const points = (geoSummary?.countries || []).filter((x) => Number.isFinite(x.lat) && Number.isFinite(x.lon));
    if (!worldMapReady) {
      return {
        grid: { left: 16, right: 16, top: 20, bottom: 22 },
        xAxis: {
          type: 'category',
          axisLabel: { color: 'rgba(148,163,184,0.8)', rotate: 20 },
          data: points.slice(0, 8).map((x) => x.country || 'Unknown'),
        },
        yAxis: {
          type: 'value',
          axisLabel: { color: 'rgba(148,163,184,0.8)' },
          splitLine: { lineStyle: { color: 'rgba(148,163,184,0.12)' } },
        },
        series: [
          {
            type: 'bar',
            data: points.slice(0, 8).map((x) => x.count),
            itemStyle: { color: '#22d3ee' },
          },
        ],
      };
    }
    return {
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'item',
        formatter: (params: { name?: string; value?: unknown }) => {
          const val = Array.isArray(params.value) ? params.value : [];
          const count = typeof val[2] === 'number' ? val[2] : 0;
          return `${params.name || '未知'}<br/>访问量：${count}`;
        },
      },
      geo: {
        map: 'world',
        roam: true,
        zoom: 1.1,
        itemStyle: {
          areaColor: 'rgba(15,23,42,0.96)',
          borderColor: 'rgba(100,116,139,0.55)',
          borderWidth: 0.6,
        },
        emphasis: {
          itemStyle: { areaColor: 'rgba(56,189,248,0.28)' },
          label: { show: false },
        },
      },
      series: [
        {
          name: '访问点位',
          type: 'effectScatter',
          coordinateSystem: 'geo',
          data: points.map((x) => ({
            name: x.country,
            value: [x.lon, x.lat, x.count],
          })),
          rippleEffect: { brushType: 'stroke', scale: 2.5 },
          itemStyle: { color: '#22d3ee' },
          symbolSize: (val: number[]) => Math.max(6, Math.min(24, (val[2] || 0) * 2)),
          encode: { value: 2 },
        },
      ],
    };
  }, [geoSummary, worldMapReady]);

  const topCountries = useMemo(() => (geoSummary?.countries || []).slice(0, 5), [geoSummary]);

  const gaugeHealth = useMemo((): EChartsOption => {
    const pct = healthOk === null ? 50 : healthOk ? 94 : 38;
    return {
      series: [
        {
          type: 'gauge',
          startAngle: 200,
          endAngle: -20,
          min: 0,
          max: 100,
          splitNumber: 5,
          radius: '88%',
          center: ['50%', '58%'],
          axisLine: {
            lineStyle: {
              width: 14,
              color: [
                [0.4, 'rgba(248,113,113,0.6)'],
                [0.65, 'rgba(251,191,36,0.7)'],
                [1, 'rgba(56,189,248,0.95)'],
              ],
            },
          },
          pointer: { length: '68%', width: 5 },
          axisTick: { distance: -14, length: 6 },
          splitLine: { distance: -18, length: 14 },
          axisLabel: { color: 'rgba(148,163,184,0.9)', fontSize: 10, distance: 14 },
          detail: {
            valueAnimation: true,
            fontSize: 22,
            fontWeight: 700,
            color: healthOk ? '#22d3ee' : '#f87171',
            formatter: '{value}',
          },
          data: [{ value: pct, name: '综合' }],
          title: { offsetCenter: [0, '78%'], color: 'rgba(148,163,184,0.85)', fontSize: 12 },
        },
      ],
    };
  }, [healthOk]);

  const gaugeRedis = useMemo((): EChartsOption => {
    const hit = cacheStats?.hit_rate;
    const v = hit === undefined || hit === null ? NaN : Math.round(hit * 100);
    return {
      series: [
        {
          type: 'gauge',
          startAngle: 200,
          endAngle: -20,
          min: 0,
          max: 100,
          radius: '88%',
          center: ['50%', '58%'],
          axisLine: {
            lineStyle: {
              width: 10,
              color: [[1, 'rgba(52,211,153,0.85)']],
            },
          },
          pointer: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          detail: {
            fontSize: 18,
            fontWeight: 700,
            color: '#34d399',
            formatter: () => (Number.isFinite(v) ? `${v}%` : '—'),
          },
          data: [{ value: Number.isFinite(v) ? v : 0 }],
          title: { offsetCenter: [0, '78%'], color: 'rgba(148,163,184,0.75)', fontSize: 11 },
        },
      ],
    };
  }, [cacheStats]);

  const formatLogTime = (t: number) => dayjs.unix(t).format('MM-DD HH:mm:ss');

  const logColumns = useMemo(
    () => [
      {
        title: '时间',
        dataIndex: 'created_at',
        key: 'created_at',
        width: 132,
        render: (t: number) => <span style={{ color: 'rgba(186,230,253,0.95)' }}>{formatLogTime(t)}</span>,
      },
      {
        title: '用户',
        dataIndex: 'username',
        key: 'username',
        width: 88,
        ellipsis: true,
      },
      {
        title: '模块',
        dataIndex: 'module',
        width: 96,
        ellipsis: true,
        render: (m: string) => m || '—',
      },
      {
        title: '操作',
        dataIndex: 'action',
        width: 72,
        render: (a: string) => <Tag color="cyan">{a || '—'}</Tag>,
      },
      {
        title: '路径',
        dataIndex: 'path',
        ellipsis: true,
      },
      {
        title: '状态',
        dataIndex: 'status',
        width: 64,
        render: (s: number) => (
          <Tag color={s >= 200 && s < 300 ? 'success' : 'error'}>{s}</Tag>
        ),
      },
      {
        title: '耗时',
        dataIndex: 'latency_ms',
        width: 64,
        render: (ms: number) => (ms != null ? `${ms}ms` : '—'),
      },
    ],
    [],
  );

  const dash = (v: number | null | undefined) => (v == null ? '—' : v);

  return (
    <PageContainer ghost title={false} childrenContentStyle={{ paddingBlockStart: 0, paddingInline: 0 }}>
      <div ref={boardRef} className="data-board">
        <div className="board-scan" aria-hidden />

        <Spin spinning={loading} size="large">
          <Row align="middle" justify="space-between" style={{ marginBottom: 20, position: 'relative', zIndex: 1 }}>
            <Col>
              <div className="board-header-title">运维态势全景 · Operations Command</div>
              <div className="board-sub">
                {userName} · 菜单权限 {menuCount} 项 · 数据每 60 秒刷新 · 无权限的接口将显示「—」
              </div>
            </Col>
            <Col>
              <Space align="center" size="middle" wrap>
                <Typography.Text style={{ color: 'rgba(56,189,248,0.85)', fontSize: 12, letterSpacing: '0.12em' }}>
                  SHIYU ADMIN DASHBOARD
                </Typography.Text>
                <Button
                  type="primary"
                  ghost
                  size="small"
                  icon={isFullscreen ? <FullscreenExitOutlined /> : <FullscreenOutlined />}
                  onClick={() => void toggleFullscreen()}
                >
                  {isFullscreen ? '退出全屏' : '全屏'}
                </Button>
              </Space>
            </Col>
          </Row>

          <Row gutter={[14, 14]} style={{ position: 'relative', zIndex: 1 }}>
            <Col xs={24} sm={12} md={6}>
              <div className="board-kpi">
                <div className="kpi-label">用户总数</div>
                <div className="kpi-value">{dash(userTotal)}</div>
                <div className="kpi-extra">系统用户 · 用户管理</div>
              </div>
            </Col>
            <Col xs={24} sm={12} md={6}>
              <div className="board-kpi">
                <div className="kpi-label">角色 / 部门</div>
                <div className="kpi-value">
                  {dash(roleTotal)}
                  <span style={{ fontSize: 18, opacity: 0.45, margin: '0 6px' }}>/</span>
                  {dash(deptCount)}
                </div>
                <div className="kpi-extra">角色管理 · 部门管理</div>
              </div>
            </Col>
            <Col xs={24} sm={12} md={6}>
              <div className="board-kpi">
                <div className="kpi-label">操作日志总条数</div>
                <div className="kpi-value">{dash(logTotal)}</div>
                <div className="kpi-extra">审计轨迹 · 操作日志</div>
              </div>
            </Col>
            <Col xs={24} sm={12} md={6}>
              <div className="board-kpi">
                <div className="kpi-label">在线会话</div>
                <div className="kpi-value">{onlineUsers.length}</div>
                <div className="kpi-extra">系统监控 · 在线用户</div>
              </div>
            </Col>
          </Row>

          <Row gutter={[14, 14]} style={{ marginTop: 14, position: 'relative', zIndex: 1 }}>
            <Col xs={24} lg={10}>
              <div className="board-panel" style={{ paddingBottom: 8 }}>
                <div className="board-chart-title">近 24 小时操作节律</div>
                <ChartShell option={lineOption} style={{ height: 280 }} />
              </div>
            </Col>
            <Col xs={24} lg={7}>
              <div className="board-panel" style={{ paddingBottom: 8 }}>
                <div className="board-chart-title">用户全球访问分布（近 24h）</div>
                <ChartShell option={geoMapOption} style={{ height: 280 }} />
                <div style={{ padding: '0 12px 10px' }}>
                  <div style={{ fontSize: 12, color: 'rgba(148,163,184,0.85)', marginBottom: 8 }}>国家访问 TOP5</div>
                  {topCountries.length === 0 ? (
                    <div style={{ fontSize: 12, color: 'rgba(148,163,184,0.72)' }}>暂无可定位访问数据</div>
                  ) : (
                    topCountries.map((item, idx) => (
                      <div
                        key={`${item.country_code}-${item.country}-${idx}`}
                        style={{
                          display: 'flex',
                          justifyContent: 'space-between',
                          fontSize: 12,
                          color: '#cbd5e1',
                          marginBottom: 4,
                        }}
                      >
                        <span>{idx + 1}. {item.country || 'Unknown'}</span>
                        <span style={{ color: '#67e8f9' }}>{item.count}</span>
                      </div>
                    ))
                  )}
                </div>
              </div>
            </Col>
            <Col xs={24} lg={7}>
              <Row gutter={[14, 14]}>
                <Col span={24}>
                  <div className="board-panel" style={{ paddingBottom: 4 }}>
                    <div className="board-chart-title">核心服务与缓存命中</div>
                    <Row>
                      <Col span={12}>
                        <ChartShell option={gaugeHealth} style={{ height: 200 }} />
                      </Col>
                      <Col span={12}>
                        <ChartShell option={gaugeRedis} style={{ height: 200 }} />
                        <div style={{ textAlign: 'center', marginTop: -8, fontSize: 11, color: 'rgba(148,163,184,0.75)' }}>
                          Redis 命中率
                        </div>
                      </Col>
                    </Row>
                  </div>
                </Col>
              </Row>
            </Col>
          </Row>

          <Row gutter={[14, 14]} style={{ marginTop: 14, position: 'relative', zIndex: 1 }}>
            <Col xs={24} lg={14}>
              <div className="board-panel">
                <div className="board-chart-title">实时操作流 · 最近 {Math.min(logsSample.length, 120)} 条</div>
                <Table<OperationLog>
                  className="board-table"
                  rowKey={(r) => `${r.id}-${r.created_at}`}
                  columns={logColumns}
                  dataSource={logsSample.slice(0, 14)}
                  pagination={false}
                  size="small"
                  scroll={{ x: 900 }}
                />
              </div>
            </Col>
            <Col xs={24} lg={5}>
              <div className="board-panel" style={{ minHeight: 320, padding: '12px 14px' }}>
                <div className="board-chart-title">数据库 · 连接池</div>
                <div style={{ fontSize: 13, color: 'rgba(226,232,240,0.9)', lineHeight: 1.9 }}>
                  <div>
                    <span style={{ color: 'rgba(148,163,184,0.85)' }}>状态 </span>
                    {dbStats?.status || '—'}
                  </div>
                  <div>
                    <span style={{ color: 'rgba(148,163,184,0.85)' }}>驱动 </span>
                    {dbStats?.driver || '—'}
                  </div>
                  <div>
                    <span style={{ color: 'rgba(148,163,184,0.85)' }}>版本 </span>
                    {dbStats?.version || '—'}
                  </div>
                  <div>
                    <span style={{ color: 'rgba(148,163,184,0.85)' }}>表数量 </span>
                    {dash(dbStats?.table_count)}
                  </div>
                  <div>
                    <span style={{ color: 'rgba(148,163,184,0.85)' }}>活跃 / 空闲连接 </span>
                    {dash(dbStats?.in_use)} / {dash(dbStats?.idle)}
                  </div>
                </div>
              </div>
            </Col>
            <Col xs={24} lg={5}>
              <div className="board-panel" style={{ minHeight: 320, padding: '12px 14px' }}>
                <div className="board-chart-title">Redis 与在线用户</div>
                <div style={{ fontSize: 12, color: 'rgba(226,232,240,0.85)', lineHeight: 1.85, marginBottom: 12 }}>
                  <div>内存 {cacheStats?.used_memory_human || '—'}</div>
                  <div>键数量 {dash(cacheStats?.db_size)}</div>
                  <div>客户端 {dash(cacheStats?.connected_clients)}</div>
                </div>
                <div style={{ maxHeight: 200, overflow: 'auto', marginBottom: 10 }}>
                  {geoSummary?.top_ips?.length ? (
                    geoSummary.top_ips.slice(0, 5).map((x) => (
                      <div key={x.ip} className="board-online-item">
                        <div style={{ fontWeight: 600, color: '#7dd3fc' }}>{x.ip}</div>
                        <div style={{ color: 'rgba(148,163,184,0.85)', fontSize: 11 }}>
                          {x.country || 'Unknown'} · {x.city || 'Unknown'} · {x.count} 次
                        </div>
                      </div>
                    ))
                  ) : (
                    <div style={{ color: 'rgba(148,163,184,0.75)', fontSize: 12 }}>暂无可定位 IP 数据</div>
                  )}
                </div>
                <div style={{ maxHeight: 140, overflow: 'auto' }}>
                  {onlineUsers.length === 0 ? (
                    <div style={{ color: 'rgba(148,163,184,0.75)', fontSize: 12 }}>暂无在线用户</div>
                  ) : (
                    onlineUsers.slice(0, 4).map((u) => (
                      <div key={u.user_code} className="board-online-item">
                        <div style={{ fontWeight: 600, color: '#7dd3fc' }}>{u.username}</div>
                        <div style={{ color: 'rgba(148,163,184,0.85)', fontSize: 11 }}>{u.ip || '—'}</div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            </Col>
          </Row>
        </Spin>
      </div>
    </PageContainer>
  );
};

export default Dashboard;
