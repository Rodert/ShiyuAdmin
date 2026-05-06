import { request } from "@umijs/max";

export type RedisValueType =
  | "all"
  | "string"
  | "list"
  | "set"
  | "zset"
  | "hash"
  | "stream";

export interface CacheDatabase {
  db: number;
  key_count: number;
}

export interface CacheKey {
  key: string;
  type: string;
  ttl: number;
}

export interface CacheKeyList {
  db: number;
  pattern: string;
  type: string;
  total: number;
  items: CacheKey[];
}

export interface CacheValue {
  db: number;
  key: string;
  type: string;
  ttl: number;
  value: any;
}

export async function getCacheDatabases() {
  return request<{
    code: number;
    data: CacheDatabase[];
    message?: string;
  }>("/api/v1/system/cache/databases", {
    method: "GET",
  });
}

export async function getCacheKeys(params: {
  db: number;
  pattern?: string;
  type?: RedisValueType;
  count?: number;
}) {
  return request<{
    code: number;
    data: CacheKeyList;
    message?: string;
  }>("/api/v1/system/cache/keys", {
    method: "GET",
    params,
  });
}

export async function getCacheValue(params: { db: number; key: string }) {
  return request<{
    code: number;
    data: CacheValue;
    message?: string;
  }>("/api/v1/system/cache/value", {
    method: "GET",
    params,
  });
}

export async function deleteCacheKey(params: { db: number; key: string }) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>("/api/v1/system/cache/key", {
    method: "DELETE",
    params,
  });
}
