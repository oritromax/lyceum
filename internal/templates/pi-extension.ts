// Lyceum's pre-tool hook for Pi, written by "lyceum init"; running it again replaces this file.
// Every tool call is sent to "lyceum hook pre-tool --harness pi", which answers allow, deny or ask.
import { spawnSync } from "node:child_process";

export default function (pi: any) {
  pi.on("tool_call", async (event: any, ctx: any) => {
    const run = spawnSync("lyceum", ["hook", "pre-tool", "--harness", "pi"], {
      input: JSON.stringify({ toolName: event.toolName, toolCallId: event.toolCallId, input: event.input }),
      encoding: "utf8",
    });
    if ((run.error as any)?.code === "ENOENT") {
      return {
        block: true,
        reason: "Lyceum blocked this tool call: lyceum is not on PATH, so its hook cannot check it. Install it with npm i -g lyceum, brew install oritromax/tap/lyceum or go install github.com/oritromax/lyceum/cmd/lyceum@latest.",
      };
    }
    let answer: any;
    try {
      answer = JSON.parse(run.stdout);
    } catch {
      return; // no answer: the call goes ahead
    }
    if (answer?.decision === "deny") {
      return { block: true, reason: answer.reason };
    }
    if (answer?.decision === "ask") {
      // Only the human, in Pi's own prompt, lets the call run; with no UI there is no one to ask.
      if (!ctx.hasUI || (await ctx.ui.confirm("Lyceum approval", answer.reason)) !== true) {
        return {
          block: true,
          reason: "Lyceum blocked this command: the owner did not confirm it in Pi's prompt. Do not retry this command or work around it.",
        };
      }
    }
  });
}
