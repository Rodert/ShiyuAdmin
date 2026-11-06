import { getNotifications } from "@/services/shiyu-api/work";
import { BellOutlined } from "@ant-design/icons";
import { history, useModel } from "@umijs/max";
import { Badge, Button, Tooltip } from "antd";
import { useEffect, useState } from "react";
export default function NotificationBell() {
  const { initialState } = useModel("@@initialState");
  const [count, setCount] = useState(0);
  const enabled =
    !!initialState?.currentUser &&
    !initialState.currentUser.requirePasswordChange;
  useEffect(() => {
    if (!enabled) return;
    let active = true;
    const load = () => {
      getNotifications({ page_size: 1 }, { skipErrorHandler: true })
        .then(({ data }) => {
          if (active) setCount(data.unread);
        })
        .catch(() => {});
    };
    load();
    const timer = setInterval(load, 60000);
    window.addEventListener("focus", load);
    window.addEventListener("shiyu-work-changed", load);
    return () => {
      active = false;
      clearInterval(timer);
      window.removeEventListener("focus", load);
      window.removeEventListener("shiyu-work-changed", load);
    };
  }, [enabled]);
  if (!enabled) return null;
  return (
    <Tooltip title="站内通知">
      <Button
        type="text"
        aria-label={`站内通知，${count} 条未读`}
        onClick={() => history.push("/work/notifications")}
      >
        <Badge count={count} size="small">
          <BellOutlined style={{ fontSize: 18 }} />
        </Badge>
      </Button>
    </Tooltip>
  );
}
