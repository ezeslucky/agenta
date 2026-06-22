import type { Metadata } from "next";
import { AgentaLanding } from "@/features/landing/components/agenta-landing";

export const metadata: Metadata = {
  title: "Homepage",
  description:
    "Agenta — open-source platform that turns coding agents into real teammates. Assign tasks, track progress, compound skills.",
  openGraph: {
    title: "Agenta — Project Management for Human + Agent Teams",
    description:
      "Manage your human + agent workforce in one place.",
    url: "/homepage",
  },
  alternates: {
    canonical: "/homepage",
  },
};

export default function HomepagePage() {
  return <AgentaLanding />;
}
