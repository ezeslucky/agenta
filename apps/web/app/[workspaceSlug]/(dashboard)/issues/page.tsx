"use client";

import { IssuesPage } from "@agenta/views/issues/components";
import { ErrorBoundary } from "@agenta/ui/components/common/error-boundary";

export default function Page() {
  return (
    <ErrorBoundary>
      <IssuesPage />
    </ErrorBoundary>
  );
}
