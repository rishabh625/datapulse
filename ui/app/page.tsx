"use client";

declare const process: {
  env: Record<string, string | undefined>;
};

import { CopilotKit } from "@copilotkit/react-core";
import { CopilotChat } from "@copilotkit/react-ui";
import React from "react";

export default function Page() {
  return (
    <CopilotKit runtimeUrl={process.env.NEXT_PUBLIC_COPILOT_RUNTIME_URL!}>
      <div style={{ height: "100vh" }}>
        <CopilotChat />
      </div>
    </CopilotKit>
  );
}