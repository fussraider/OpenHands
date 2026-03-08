import os
import re

def get_logs_from_dir(directory, regex_pattern):
    logs = []
    for root, _, files in os.walk(directory):
        for file in files:
            if not file.endswith('.py') and not file.endswith('.go'):
                continue
            filepath = os.path.join(root, file)
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.finditer(regex_pattern, content)
                for match in matches:
                    log_content = match.group(1).strip()
                    # Strip quotes and f-string prefix
                    log_content = re.sub(r"^[furb]*['\"](.*)['\"]$", r"\1", log_content)

                    # Clean up formatting strings
                    log_content = re.sub(r'\{[^}]+\}', '{}', log_content)
                    log_content = re.sub(r'%[sdvTq+]', '{}', log_content)

                    logs.append({
                        'file': filepath,
                        'raw': match.group(0),
                        'cleaned': log_content.strip()
                    })
    return logs

def main():
    python_pattern = r"(?:logger|openhands_logger)\.debug\(\s*([furb]*['\"][^'\"]*['\"])"
    go_pattern = r"slog\.Debug\(\s*\"([^\"]+)\""

    python_logs = get_logs_from_dir("openhands", python_pattern)
    go_logs = get_logs_from_dir("server", go_pattern)

    go_cleaned = [g['cleaned'].lower() for g in go_logs]

    missing = []
    for p in python_logs:
        # Check if a roughly similar string exists in Go logs
        found = False
        p_text = p['cleaned'].lower()
        if len(p_text) < 5:
            continue

        for g in go_cleaned:
            if g in p_text or p_text in g or (len(g) > 10 and len(p_text) > 10 and g[:10] == p_text[:10]):
                found = True
                break

        if not found:
            missing.append(p)

    with open('missing_debug_logs_report.md', 'w') as f:
        f.write("# Missing Debug Logs Report (Python vs Go)\n\n")
        f.write(f"Total Python debug logs parsed: {len(python_logs)}\n")
        f.write(f"Total Go debug logs parsed: {len(go_logs)}\n")
        f.write(f"Missing in Go: {len(missing)}\n\n")
        f.write("## Unported Logs\n\n")
        f.write("| Python File | Log Statement |\n")
        f.write("| --- | --- |\n")
        for m in sorted(missing, key=lambda x: x['file']):
            f.write(f"| `{m['file']}` | `{m['cleaned']}` |\n")

if __name__ == "__main__":
    main()
