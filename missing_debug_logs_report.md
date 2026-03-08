# Missing Debug Logs Report (Python vs Go)

Total Python debug logs parsed: 340
Total Go debug logs parsed: 63
Missing in Go: 253

## Unported Logs

| Python File | Log Statement |
| --- | --- |
| `openhands/agenthub/codeact_agent/codeact_agent.py` | `Using condenser: {}` |
| `openhands/agenthub/codeact_agent/codeact_agent.py` | `Processing {} events from a total of {} events` |
| `openhands/agenthub/codeact_agent/codeact_agent.py` | `Actions after response_to_actions: {}` |
| `openhands/agenthub/loc_agent/loc_agent.py` | `TOOLS loaded for LocAgent: {` |
| `openhands/agenthub/readonly_agent/readonly_agent.py` | `TOOLS loaded for ReadOnlyAgent: {` |
| `openhands/app_server/app_conversation/app_conversation_service_base.py` | `Loading skills for V1 conversation via agent-server` |
| `openhands/app_server/app_conversation/live_status_app_conversation_service.py` | `Added custom SSE server: {} for {}` |
| `openhands/app_server/app_conversation/live_status_app_conversation_service.py` | `Added custom SHTTP server: {} for {}` |
| `openhands/app_server/app_conversation/live_status_app_conversation_service.py` | `Added custom STDIO server: {}` |
| `openhands/app_server/app_conversation/skill_loader.py` | `org-level skill directory {} not found: {}` |
| `openhands/app_server/app_conversation/sql_app_conversation_info_service.py` | `No agent metrics found in stats for conversation {}` |
| `openhands/app_server/app_conversation/sql_app_conversation_info_service.py` | `Conversation {} not found or not accessible, skipping statistics update` |
| `openhands/app_server/sandbox/docker_sandbox_service.py` | `Sandbox server not yet available (still starting):` |
| `openhands/app_server/sandbox/docker_sandbox_spec_service.py` | `Checking Docker Image: {}` |
| `openhands/app_server/sandbox/remote_sandbox_service.py` | `Matched {} Runtimes with {} Conversations.` |
| `openhands/app_server/sandbox/remote_sandbox_service.py` | `Started Refreshing Conversation {}` |
| `openhands/app_server/sandbox/remote_sandbox_service.py` | `Finished Refreshing Conversation {}` |
| `openhands/app_server/sandbox/sandbox_service.py` | `Agent server health check failed for sandbox {}` |
| `openhands/controller/agent_controller.py` | `Original security risk for {}: {})` |
| `openhands/controller/agent_controller.py` | `[Security Analyzer: {}] Override security risk for action {}: {}` |
| `openhands/controller/agent_controller.py` | `No security analyzer configured, setting UNKNOWN risk for action: {}` |
| `openhands/controller/agent_controller.py` | `System message: {}` |
| `openhands/controller/agent_controller.py` | `[non-CLI mode] Detected HIGH security risk in action: {}. Ask for confirmation` |
| `openhands/controller/state/state.py` | `Saving state to session {}:{}` |
| `openhands/core/config/sandbox_config.py` | `SandboxConfig user_id default: {}` |
| `openhands/core/config/utils.py` | `Default model routing configuration loaded from config toml and assigned to default agent` |
| `openhands/core/config/utils.py` | `Default condenser configuration loaded from config toml and assigned to default agent` |
| `openhands/core/config/utils.py` | `Default LLM summarizing condenser assigned to default agent (no condenser in config)` |
| `openhands/core/config/utils.py` | `No explicit /workspace mount found in SANDBOX_VOLUMES.` |
| `openhands/core/config/utils.py` | `Automatically disabled Jupyter plugin and browsing for all agents` |
| `openhands/core/config/utils.py` | `Loading agent config from {}` |
| `openhands/core/config/utils.py` | `Loading from toml failed for {}` |
| `openhands/core/config/utils.py` | `Loading llm config` |
| `openhands/core/config/utils.py` | `Config file not found: {}` |
| `openhands/core/config/utils.py` | `LLM config` |
| `openhands/core/config/utils.py` | `Loading condenser config [{}] from {}` |
| `openhands/core/config/utils.py` | `Condenser [{}] requires LLM config [{}]. Loading it...` |
| `openhands/core/config/utils.py` | `Loading model routing config [` |
| `openhands/core/config/utils.py` | `CLI specified LLM config: {}` |
| `openhands/core/config/utils.py` | `Trying to load LLM config` |
| `openhands/core/config/utils.py` | `Using LLM config` |
| `openhands/core/config/utils.py` | `Set LLM config from CLI parameter: {}` |
| `openhands/core/logger.py` | `DEBUG mode enabled.` |
| `openhands/core/logger.py` | `Logging initialized` |
| `openhands/core/logger.py` | `Logging to file in: {}` |
| `openhands/core/logger.py` | `Logging to {}` |
| `openhands/core/main.py` | `Agent Controller Initialized: Running agent {}, model` |
| `openhands/core/main.py` | `Stopping agent controller...` |
| `openhands/core/main.py` | `Stopping EventStream...` |
| `openhands/core/main.py` | `Closing runtime...` |
| `openhands/core/setup.py` | `Runtime created with plugins: {}` |
| `openhands/core/setup.py` | `Selected repository {}.` |
| `openhands/core/setup.py` | `Trying to restore agent state from session {} if available` |
| `openhands/core/setup.py` | `Cannot restore agent state: {}` |
| `openhands/events/observation/commands.py` | `Truncated large command output: {} -> {} chars` |
| `openhands/events/utils.py` | `Event {} has no ID` |
| `openhands/events/utils.py` | `Observation {} has no cause` |
| `openhands/events/utils.py` | `Observation {} has no cause` |
| `openhands/integrations/forgejo/service/features.py` | `Forgejo microagent scan warning for {}: {}` |
| `openhands/integrations/provider.py` | `No microagents found on {} for {}, trying other providers` |
| `openhands/integrations/provider.py` | `No content found on {} for {}/{}, trying other providers` |
| `openhands/integrations/provider.py` | `File not found on {} for {}/{}, trying other providers` |
| `openhands/integrations/provider.py` | `[Azure DevOps] Original domain: {}` |
| `openhands/integrations/provider.py` | `[Azure DevOps] Token available: {},` |
| `openhands/integrations/provider.py` | `[Azure DevOps] Cleaned domain: {}` |
| `openhands/integrations/provider.py` | `[Azure DevOps] Repository parts: {} (length: {})` |
| `openhands/integrations/provider.py` | `[Azure DevOps] URL-encoded parts - org: {}, project: {}, repo: {}` |
| `openhands/integrations/service_types.py` | `No .cursorrules file found in {}` |
| `openhands/llm/async_llm.py` | `LLM request cancelled by user.` |
| `openhands/llm/llm.py` | `LLM: caching prompt enabled` |
| `openhands/llm/llm.py` | `Gemini model {} with reasoning_effort {}` |
| `openhands/llm/llm.py` | `Gemini model {} with reasoning_effort {} mapped to thinking {kwargs.get(` |
| `openhands/llm/llm.py` | `Model info: {json.dumps({` |
| `openhands/llm/llm.py` | `Setting top_p to 0.9 for Hugging Face model: {}` |
| `openhands/llm/llm.py` | `Using context window: {}` |
| `openhands/llm/llm.py` | `Using custom cost per token: {}` |
| `openhands/llm/llm.py` | `Got response_cost from response: {}` |
| `openhands/llm/llm.py` | `Using fallback model name {} to get cost: {}` |
| `openhands/llm/router/base.py` | `RouterLLM routing to {} ({})` |
| `openhands/llm/streaming_llm.py` | `LLM request cancelled by user.` |
| `openhands/mcp/utils.py` | `Set SHTTP server timeout to {}s` |
| `openhands/mcp/utils.py` | `Creating MCP clients with config: {}` |
| `openhands/mcp/utils.py` | `No MCP clients were successfully connected` |
| `openhands/memory/condenser/impl/conversation_window_condenser.py` | `Removed {} dangling observation(s)` |
| `openhands/memory/conversation_memory.py` | `Visual browsing: {}` |
| `openhands/memory/conversation_memory.py` | `IPython observation has image URLs but none are valid` |
| `openhands/memory/conversation_memory.py` | `Adding {} for browsing` |
| `openhands/memory/conversation_memory.py` | `Vision enabled for browsing, but no valid image available` |
| `openhands/memory/conversation_memory.py` | `[ConversationMemory] No SystemMessageAction found in events.` |
| `openhands/memory/conversation_memory.py` | `The user MessageAction at index 1 does not match the provided initial_user_action.` |
| `openhands/memory/memory.py` | `Workspace context recall` |
| `openhands/memory/memory.py` | `Found MCP tools in repo microagent {}: {}` |
| `openhands/microagent/microagent.py` | `This microagent requires user input: {}` |
| `openhands/microagent/microagent.py` | `Loading agents from {}` |
| `openhands/microagent/microagent.py` | `Loaded {} microagents:` |
| `openhands/runtime/action_execution_server.py` | `Browser initialized asynchronously` |
| `openhands/runtime/action_execution_server.py` | `Browser is ready` |
| `openhands/runtime/action_execution_server.py` | `Bash session initialized` |
| `openhands/runtime/action_execution_server.py` | `Browser initialization started in background` |
| `openhands/runtime/action_execution_server.py` | `All plugins initialized` |
| `openhands/runtime/action_execution_server.py` | `AgentSkills initialized: {}` |
| `openhands/runtime/action_execution_server.py` | `Runtime client initialized.` |
| `openhands/runtime/action_execution_server.py` | `Init command outputs (exit code: {}): {}` |
| `openhands/runtime/action_execution_server.py` | `Bash init commands completed` |
| `openhands/runtime/action_execution_server.py` | `{} != {} -> reset Jupyter PWD` |
| `openhands/runtime/action_execution_server.py` | `Changed working directory in IPython to: {}. Output: {}` |
| `openhands/runtime/action_execution_server.py` | `Uploaded file {} and extracted to {}` |
| `openhands/runtime/action_execution_server.py` | `Uploaded file {} to {}` |
| `openhands/runtime/action_execution_server.py` | `Downloading files` |
| `openhands/runtime/action_execution_server.py` | `Starting action execution API on port {}` |
| `openhands/runtime/base.py` | `Security analyzer {} initialized for runtime {}` |
| `openhands/runtime/base.py` | `Adding env vars: {}` |
| `openhands/runtime/base.py` | `Added env vars to IPython` |
| `openhands/runtime/base.py` | `Adding env vars to PowerShell` |
| `openhands/runtime/base.py` | `Added env vars to PowerShell session: {}` |
| `openhands/runtime/base.py` | `Adding env vars to bash` |
| `openhands/runtime/base.py` | `Adding env var to .bashrc: {}` |
| `openhands/runtime/base.py` | `No repository selected. Initializing a new git repository in the workspace.` |
| `openhands/runtime/base.py` | `Skipping git configuration for CLI runtime - using user` |
| `openhands/runtime/browser/browser_env.py` | `Starting browser env...` |
| `openhands/runtime/browser/browser_env.py` | `Browsing goal: {}` |
| `openhands/runtime/browser/browser_env.py` | `SHUTDOWN recv, shutting down browser env...` |
| `openhands/runtime/browser/browser_env.py` | `Browser env process interrupted by user.` |
| `openhands/runtime/browser/browser_env.py` | `Browser env is not alive. Response ID: {}` |
| `openhands/runtime/builder/docker.py` | `Checking, if image exists locally:\n{}` |
| `openhands/runtime/builder/docker.py` | `Image found locally.` |
| `openhands/runtime/builder/docker.py` | `Image {} {colorize(` |
| `openhands/runtime/builder/docker.py` | `Image not found locally. Trying to pull it, please wait...` |
| `openhands/runtime/builder/docker.py` | `Image pulled` |
| `openhands/runtime/builder/docker.py` | `Layer {}: {layers[layer_id][` |
| `openhands/runtime/builder/docker.py` | `Removed old cache file: {}` |
| `openhands/runtime/builder/docker.py` | `Created cache directory: {}` |
| `openhands/runtime/builder/docker.py` | `Cache directory {} is usable` |
| `openhands/runtime/builder/remote.py` | `Image {} exists.` |
| `openhands/runtime/builder/remote.py` | `Image {} does not exist.` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `Details: {}` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[CLIRuntime] Set os.environ[` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_safe_terminate_process] Original PID to act on: {}` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_safe_terminate_process] Attempting to {} for PID {} (PGID: {}) with {}.` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_safe_terminate_process] Successfully sent signal {} to PGID {} (original PID: {}).` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_safe_terminate_process] Fallback: Terminated {} (PID: {}).` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_safe_terminate_process] Fallback: Terminated {} (PID: {}).` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_execute_shell_command] PID of bash -c: {} for command:` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `[_execute_shell_command] Complete output for` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `Running command in CLIRuntime:` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `Skipping recursive copy: source and target are identical.` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `Skipping copy as source and destination are the same: {}` |
| `openhands/runtime/impl/cli/cli_runtime.py` | `PowerShell session closed successfully.` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Mount dir (sandbox.volumes): {} to {} with mode: {}` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Mount dir (legacy): {} with mode: {}` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Mount dir is not set, will not mount the workspace directory to the container` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Released host port lock for port {}` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Released VSCode port lock for port {}` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Released app port lock for port {self._app_ports[i] if i < len(self._app_ports) else` |
| `openhands/runtime/impl/docker/docker_runtime.py` | `Port {} is in use by Docker, trying again` |
| `openhands/runtime/impl/local/local_runtime.py` | `Checking dependencies: Jupyter` |
| `openhands/runtime/impl/local/local_runtime.py` | `Jupyter output: {}` |
| `openhands/runtime/impl/local/local_runtime.py` | `Checking dependencies: libtmux` |
| `openhands/runtime/impl/local/local_runtime.py` | `Checking dependencies: browser` |
| `openhands/runtime/impl/local/local_runtime.py` | `runtime_url is {}` |
| `openhands/runtime/impl/local/local_runtime.py` | `_create_url url is {}` |
| `openhands/runtime/impl/local/local_runtime.py` | `Updated PATH for subprocesses: {env[` |
| `openhands/runtime/impl/local/local_runtime.py` | `Warm server not ready yet: {}` |
| `openhands/runtime/impl/remote/remote_runtime.py` | `RemoteRuntime.init user_id {}` |
| `openhands/runtime/plugins/jupyter/__init__.py` | `Jupyter launch command (Windows): {}` |
| `openhands/runtime/plugins/jupyter/__init__.py` | `Jupyter kernel gateway started at port {}. Output: {}` |
| `openhands/runtime/plugins/jupyter/__init__.py` | `Jupyter launch command: {}` |
| `openhands/runtime/plugins/jupyter/__init__.py` | `Jupyter kernel gateway started at port {}. Output: {}` |
| `openhands/runtime/plugins/vscode/__init__.py` | `VSCode server started at port {}. Output: {}` |
| `openhands/runtime/plugins/vscode/__init__.py` | `VSCode settings copied to {}` |
| `openhands/runtime/utils/bash.py` | `BASH PARSING between: {}` |
| `openhands/runtime/utils/bash.py` | `BASH PARSING command: {}` |
| `openhands/runtime/utils/bash.py` | `BASH PARSING remaining: {}` |
| `openhands/runtime/utils/bash.py` | `BASH PARSING result[-1] += remaining: {}` |
| `openhands/runtime/utils/bash.py` | `BASH PARSING result.append(remaining): {}` |
| `openhands/runtime/utils/bash.py` | `pane: {}; history_limit: {}` |
| `openhands/runtime/utils/bash.py` | `Bash session initialized with work dir: {}` |
| `openhands/runtime/utils/bash.py` | `directory_changed: {}; {}; {}` |
| `openhands/runtime/utils/bash.py` | `COMMAND OUTPUT: {}` |
| `openhands/runtime/utils/bash.py` | `COMBINED OUTPUT: {}` |
| `openhands/runtime/utils/bash.py` | `RECEIVED ACTION: {}` |
| `openhands/runtime/utils/bash.py` | `Initial PS1 count: {}` |
| `openhands/runtime/utils/bash.py` | `PREVIOUS COMMAND OUTPUT: {}` |
| `openhands/runtime/utils/bash.py` | `SENDING INPUT TO RUNNING PROCESS: {}` |
| `openhands/runtime/utils/bash.py` | `SENDING COMMAND: {}` |
| `openhands/runtime/utils/bash.py` | `GETTING PANE CONTENT at {}` |
| `openhands/runtime/utils/bash.py` | `PANE CONTENT GOT after {} seconds` |
| `openhands/runtime/utils/bash.py` | `PANE_CONTENT: {}` |
| `openhands/runtime/utils/bash.py` | `BEGIN OF PANE CONTENT: {}` |
| `openhands/runtime/utils/bash.py` | `END OF PANE CONTENT: {}` |
| `openhands/runtime/utils/bash.py` | `CONTENT UPDATED DETECTED at {}` |
| `openhands/runtime/utils/bash.py` | `CHECKING NO CHANGE TIMEOUT ({}s): elapsed {}. Action blocking: {}` |
| `openhands/runtime/utils/bash.py` | `CHECKING HARD TIMEOUT ({}s): elapsed {}` |
| `openhands/runtime/utils/bash.py` | `Hard timeout triggered.` |
| `openhands/runtime/utils/bash.py` | `SLEEPING for {} seconds for next poll` |
| `openhands/runtime/utils/command.py` | `app_config {}` |
| `openhands/runtime/utils/command.py` | `sandbox_config {}` |
| `openhands/runtime/utils/command.py` | `RUNTIME_USERNAME {}, RUNTIME_UID {}` |
| `openhands/runtime/utils/command.py` | `override_username {}, override_user_id {}` |
| `openhands/runtime/utils/command.py` | `username {}, user_id {}` |
| `openhands/runtime/utils/command.py` | `get_action_execution_server_startup_command: {}` |
| `openhands/runtime/utils/edit.py` | `It is not recommended to cache draft editor LLM prompts as it may incur high costs for the same prompt.` |
| `openhands/runtime/utils/edit.py` | `[Draft edit functionality] enabled with LLM: {}` |
| `openhands/runtime/utils/edit.py` | `Agent attempted to edit a file that does not exist. Creating the file. Error msg: {}` |
| `openhands/runtime/utils/port_lock.py` | `Acquired lock for port {}` |
| `openhands/runtime/utils/port_lock.py` | `Acquired lock for port {}` |
| `openhands/runtime/utils/port_lock.py` | `Released lock for port {}` |
| `openhands/runtime/utils/port_lock.py` | `Found and locked available port {}` |
| `openhands/runtime/utils/port_lock.py` | `Found and locked available port {}` |
| `openhands/runtime/utils/port_lock.py` | `Cleaned up stale lock file: {}` |
| `openhands/runtime/utils/runtime_build.py` | `The provided image [{}] is already a valid runtime image.\n` |
| `openhands/runtime/utils/runtime_build.py` | `Force rebuild: [{}:{}] from scratch.` |
| `openhands/runtime/utils/runtime_build.py` | `Reusing Image [{}]` |
| `openhands/runtime/utils/runtime_build.py` | `Build [{}] from lock image [{}]` |
| `openhands/runtime/utils/runtime_build.py` | `Build [{}] from scratch` |
| `openhands/runtime/utils/runtime_build.py` | `Building source distribution using project root: {}` |
| `openhands/runtime/utils/runtime_build.py` | `Copying the source code and generating the Dockerfile in the build folder: {}` |
| `openhands/runtime/utils/runtime_build.py` | `Runtime image repo: {} and runtime image tag: {}` |
| `openhands/runtime/utils/runtime_build.py` | `Build folder [{}] is ready: {}` |
| `openhands/runtime/utils/runtime_build.py` | ``config.sh` is updated with the image repo[{}] and tags [{}, {}]` |
| `openhands/runtime/utils/runtime_build.py` | `Dockerfile, source code and config.sh are ready in {}` |
| `openhands/runtime/utils/runtime_build.py` | `Building image in a temporary folder` |
| `openhands/runtime/utils/runtime_build.py` | `\nBuilt image: {}\n` |
| `openhands/runtime/utils/runtime_init.py` | `Running on Windows, skipping Unix-specific user setup` |
| `openhands/runtime/utils/runtime_init.py` | `Client working directory: {}` |
| `openhands/runtime/utils/runtime_init.py` | `Created working directory: {}` |
| `openhands/runtime/utils/runtime_init.py` | `Attempting to create user `{}` with UID {}.` |
| `openhands/runtime/utils/runtime_init.py` | `User `{}` already has the provided UID {}. Skipping user setup.` |
| `openhands/runtime/utils/runtime_init.py` | `User `{}` does not exist. Proceeding with user creation.` |
| `openhands/runtime/utils/runtime_init.py` | `Added sudoer successfully. Output: [{}]` |
| `openhands/runtime/utils/runtime_init.py` | `Added user `{}` successfully with UID {}. Output: [{}]` |
| `openhands/runtime/utils/runtime_init.py` | `Client working directory: {}` |
| `openhands/runtime/utils/runtime_init.py` | `Created working directory. Output: [{}]` |
| `openhands/runtime/utils/windows_bash.py` | `Imported clr module from: {}` |
| `openhands/runtime/utils/windows_bash.py` | `Loaded PowerShell SDK (Desktop): {}` |
| `openhands/runtime/utils/windows_bash.py` | `Confirmed runspace CWD is {}` |
| `openhands/runtime/utils/windows_bash.py` | `Running PS command:` |
| `openhands/security/grayswan/analyzer.py` | `Event stream set for GraySwanAnalyzer` |
| `openhands/security/grayswan/analyzer.py` | `Calling security_risk on GraySwanAnalyzer for action: {}` |
| `openhands/security/grayswan/analyzer.py` | `Converted {} events into {} OpenAI messages for GraySwan analysis` |
| `openhands/security/invariant/analyzer.py` | `Calling security_risk on InvariantAnalyzer` |
| `openhands/server/listen_socket.py` | `Invalid latest_event_id value: {}, defaulting to -1` |
| `openhands/server/routes/settings.py` | `No api_base found in litellm for model: {}` |
| `openhands/storage/local.py` | `Local path does not exist: {}` |
| `openhands/storage/local.py` | `Removed local file: {}` |
| `openhands/storage/local.py` | `Removed local directory: {}` |
| `openhands/storage/memory.py` | `Cleared in-memory file store: {}` |
| `openhands/utils/chunk_localizer.py` | `Language {} not supported. Falling back to raw string.` |
| `openhands/utils/http_session.py` | `HttpSession:request called with args {} and kwargs {}` |
| `openhands/utils/shutdown_listener.py` | `shutdown_signal:{}` |
| `openhands/utils/shutdown_listener.py` | `_register_signal_handlers` |
| `openhands/utils/shutdown_listener.py` | `_register_signal_handlers:main_thread` |
| `openhands/utils/shutdown_listener.py` | `_register_signal_handlers:not_main_thread` |
