// Lyceum's pre-tool hook for OpenCode, written by "lyceum init"; running it again replaces this file.
// Every tool call is sent to "lyceum hook pre-tool --harness opencode", with the project directory as cwd,
// and lyceum answers allow, deny or ask. On ask, the permission check of the same call asks the human.
import { spawnSync } from "node:child_process";

export default {
  id: "lyceum",
  async setup(ctx) {
    const decide = (payload) => {
      const run = spawnSync("lyceum", ["hook", "pre-tool", "--harness", "opencode"], {
        input: JSON.stringify({ ...payload, cwd: ctx.location.directory }),
        encoding: "utf8",
      });
      if (run.error?.code === "ENOENT") {
        return {
          decision: "deny",
          reason: "Lyceum blocked this tool call: lyceum is not on PATH, so its hook cannot check it. Install it with npm i -g lyceum, brew install oritromax/tap/lyceum or go install github.com/oritromax/lyceum/cmd/lyceum@latest.",
        };
      }
      try {
        return JSON.parse(run.stdout);
      } catch {
        return undefined; // no answer: the call goes ahead
      }
    };
    await ctx.tool.hook("execute.before", (event) => {
      const answer = decide({ ...event, tool: event.tool, input: event.input });
      if (answer?.decision === "deny") {
        throw new Error(answer.reason);
      }
    });
    await ctx.permission.hook("evaluate", (event) => {
      if (event.action !== "shell") {
        return;
      }
      const answer = decide({ tool: "shell", input: { command: event.resources[0] } });
      if (answer?.decision === "ask" || answer?.decision === "deny") {
        event.effect = answer.decision;
        event.message = answer.reason;
      }
    });
  },
};
