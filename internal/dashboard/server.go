package dashboard  

import (
	"html/template"
	"net/http"
	"strconv"
	"vexora-studio/internal/database"
)

// We embed the HTML in the binary so you only need one executable
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vexora Studio</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <style>
        .markdown-preview h1 { font-size: 1.5em; font-weight: bold; margin-bottom: 0.5em; }
        .markdown-preview h2 { font-size: 1.25em; font-weight: bold; margin-top: 1em; margin-bottom: 0.5em; }
        .markdown-preview ul { list-style-type: disc; margin-left: 1.5em; }
        .markdown-preview p { margin-bottom: 1em; }
        .markdown-preview pre { background: #f3f4f6; padding: 1em; border-radius: 0.5em; overflow-x: auto; }
    </style>
</head>
<body class="bg-gray-100 h-screen flex flex-col">
    <header class="bg-gray-900 text-white p-4 flex justify-between items-center shadow-lg">
        <h1 class="text-xl font-bold tracking-wider">VEXORA STUDIO</h1>
        <div class="text-sm text-gray-400">Status: <span class="text-green-400">Online</span></div>
    </header>

    <div class="flex flex-1 overflow-hidden">
        <aside class="w-1/4 bg-white border-r border-gray-200 overflow-y-auto">
            <div class="p-4 border-b bg-gray-50 font-semibold text-gray-700">Queue ({{len .Jobs}})</div>
            <ul>
                {{range .Jobs}}
                <li class="border-b hover:bg-blue-50 transition cursor-pointer p-4 group {{if eq .ID $.SelectedID}}bg-blue-100{{end}}">
                    <a href="/dashboard?id={{.ID}}" class="block">
                        <div class="flex justify-between items-start mb-1">
                            <span class="font-bold text-gray-800">{{.ProjectName}}</span>
                            <span class="text-xs px-2 py-0.5 rounded-full {{if eq .Status "WAITING_APPROVAL"}}bg-yellow-200 text-yellow-800{{else}}bg-gray-200{{end}}">
                                {{.Status}}
                            </span>
                        </div>
                        <div class="text-xs text-gray-500 truncate">{{.GeneratedSubject}}</div>
                        <div class="text-xs text-gray-400 mt-2">{{.CreatedAt}}</div>
                    </a>
                </li>
                {{end}}
            </ul>
        </aside>

        <main class="flex-1 flex flex-col bg-gray-50 overflow-y-auto p-8">
            {{if .SelectedJob}}
                <div class="bg-white rounded-lg shadow-xl flex-1 flex flex-col overflow-hidden">
                    <div class="p-4 border-b flex justify-between items-center bg-gray-50">
                        <div>
                            <h2 class="text-lg font-bold">{{.SelectedJob.GeneratedSubject}}</h2>
                            <p class="text-xs text-gray-500">ID: {{.SelectedJob.ID}} • Type: Newsletter</p>
                        </div>
                        <div class="space-x-2">
                            <button onclick="rejectJob({{.SelectedJob.ID}})" class="px-4 py-2 bg-red-100 text-red-700 rounded hover:bg-red-200 font-medium">Reject</button>
                            <button onclick="approveJob({{.SelectedJob.ID}})" class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 font-medium shadow-md">Approve & Publish</button>
                        </div>
                    </div>

                    <div class="flex flex-1 overflow-hidden">
                        <textarea id="editor" class="w-1/2 p-6 font-mono text-sm bg-gray-900 text-gray-100 resize-none focus:outline-none" oninput="updatePreview()">{{.SelectedJob.GeneratedContent}}</textarea>
                        
                        <div id="preview" class="w-1/2 p-6 overflow-y-auto markdown-preview prose max-w-none"></div>
                    </div>
                </div>
            {{else}}
                <div class="flex items-center justify-center h-full text-gray-400">
                    <div class="text-center">
                        <div class="text-6xl mb-4">📭</div>
                        <p>Select a job from the queue to review.</p>
                    </div>
                </div>
            {{end}}
        </main>
    </div>

    <script>
        // Simple Markdown Renderer
        function updatePreview() {
            const raw = document.getElementById('editor').value;
            // Remove the JSON wrapping if it exists (fixes the artifact from your logs)
            // Avoid using literal backticks in the Go raw string by using Unicode escapes for the backtick character
            let clean = raw;
            if (raw.startsWith('\u0060\u0060\u0060markdown')) {
                clean = raw.replace(/^\u0060\u0060\u0060markdown\n/, '').replace(/\n\u0060\u0060\u0060$/, '');
            }
            document.getElementById('preview').innerHTML = marked.parse(clean);
        }

        // Run once on load
        if(document.getElementById('editor')) {
            updatePreview();
        }

        // API Calls
        async function approveJob(id) {
            if(!confirm("Ready to approve?")) return;
            // You'd need to fetch the token from the UI logic or update API to be cookie-based
            // For local dev, we can simplify or assume a mock token for now
            alert("Approval Logic triggered for " + id);
            // fetch('/approve?id=' + id ...);
        }
        
        async function rejectJob(id) {
            if(!confirm("Reject and retry?")) return;
            alert("Reject Logic triggered for " + id);
        }
    </script>
</body>
</html>
`

type PageData struct {
	Jobs        []database.QueueItem // We'll need to update QueueItem to have all fields
	SelectedJob *database.QueueItem
	SelectedID  int64
}

func StartDashboard(port string) {
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		// 1. Fetch All "WAITING_APPROVAL" Jobs
		jobs, _ := database.GetJobsByStatus("WAITING_APPROVAL")
		
		// 2. Determine Selected Job
		var selected *database.QueueItem
		idStr := r.URL.Query().Get("id")
		selectedID, _ := strconv.ParseInt(idStr, 10, 64)

		if selectedID != 0 {
			// Find it in the list (or query DB directly)
			for i := range jobs {
				if jobs[i].ID == selectedID {
					selected = &jobs[i]
					break
				}
			}
		} else if len(jobs) > 0 {
			// Default to first
			selected = &jobs[0]
			selectedID = selected.ID
		}

		// 3. Render
		tmpl, _ := template.New("dash").Parse(htmlTemplate)
		tmpl.Execute(w, PageData{
			Jobs:        jobs,
			SelectedJob: selected,
			SelectedID:  selectedID,
		})
	})

	go http.ListenAndServe(port, nil)
}