"use client";

import { DashboardLayout } from "@agenta/views/layout";
import { AgentaIcon } from "@agenta/ui/components/common/agenta-icon";
import { SearchCommand, SearchTrigger } from "@agenta/views/search";
import { ChatFab, ChatWindow } from "@agenta/views/chat";
import { WebNotificationBridge } from "@/components/web-notification-bridge";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <DashboardLayout
      loadingIndicator={<AgentaIcon className="size-6" />}
      searchSlot={<SearchTrigger />}
      extra={
        <>
          <SearchCommand />
          <ChatWindow />
          <ChatFab />
          <WebNotificationBridge />
        </>
      }
    >
      {children}
    </DashboardLayout>
  );
}
