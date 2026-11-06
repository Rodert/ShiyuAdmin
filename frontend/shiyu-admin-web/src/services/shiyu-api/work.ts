import { request } from "@umijs/max";
import type { MediaFile } from "./media";
const base = "/api/v1/system/work";
type Result<T> = { code: number; data: T };
export type WorkStatus = "pending" | "approved" | "rejected" | "withdrawn";
export type WorkItem = {
  id: string;
  title: string;
  content: string;
  applicant_code: string;
  applicant_name: string;
  approver_code: string;
  approver_name: string;
  status: WorkStatus;
  created_at: string;
};
export type WorkDetail = WorkItem & {
  attachments: { file_code: string; name: string }[];
  events: {
    id: number;
    actor_name: string;
    action: string;
    comment: string;
    created_at: string;
  }[];
};
export type Notice = {
  id: number;
  request_id: string;
  title: string;
  read_at: string | null;
  created_at: string;
};
export type Person = { user_code: string; username: string; nickname: string };
export const getPeople = (keyword = "") =>
  request<Result<Person[]>>(`${base}/people`, { params: { keyword } });
export const getRequests = (params: {
  view: string;
  page?: number;
  page_size?: number;
}) =>
  request<Result<{ items: WorkItem[]; total: number }>>(`${base}/requests`, {
    params,
  });
export const getWorkDetail = (id: string) =>
  request<Result<WorkDetail>>(`${base}/requests/${encodeURIComponent(id)}`);
export const createRequest = (data: {
  title: string;
  content: string;
  approver_code: string;
  file_codes: string[];
}) => request<Result<WorkItem>>(`${base}/requests`, { method: "POST", data });
export const actOnRequest = (id: string, action: string, comment: string) =>
  request(`${base}/requests/${encodeURIComponent(id)}/actions`, {
    method: "POST",
    data: { action, comment },
  });
export const getNotifications = (
  params: { page?: number; page_size?: number; unread?: boolean } = {},
  options: { skipErrorHandler?: boolean } = {}
) =>
  request<Result<{ items: Notice[]; total: number; unread: number }>>(
    `${base}/notifications`,
    { params, ...options }
  );
export const readNotification = (id: number) =>
  request(`${base}/notifications/${id}/read`, { method: "PUT" });
export const readAllNotifications = () =>
  request(`${base}/notifications/read-all`, { method: "PUT" });
export const notifyWorkChanged = () =>
  window.dispatchEvent(new Event("shiyu-work-changed"));
export async function uploadAttachment(file: File) {
  const data = new FormData();
  data.append("file", file);
  return request<Result<MediaFile>>(`${base}/attachments/upload`, {
    method: "POST",
    data,
    requestType: "form",
  });
}
export async function downloadAttachment(
  id: string,
  code: string,
  name: string
) {
  const token = localStorage.getItem("shiyu_token");
  const result = await fetch(
    `${base}/requests/${encodeURIComponent(
      id
    )}/attachments/${encodeURIComponent(code)}`,
    { headers: { Authorization: `Bearer ${token || ""}` } }
  );
  if (!result.ok) throw new Error("附件下载失败，请刷新后重试");
  const url = URL.createObjectURL(await result.blob());
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  a.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
