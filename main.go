package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FLV 播放器</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/flv.js/1.6.2/flv.min.js"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .input-group {
            margin-bottom: 20px;
        }
        textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
            margin-bottom: 10px;
            min-height: 100px;
            font-family: monospace;
            resize: vertical;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #45a049;
        }
        video {
            width: 100%;
            background-color: #000;
            border-radius: 4px;
        }
        .player-container {
            margin-top: 20px;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
            gap: 20px;
        }
        .player-wrapper {
            background-color: #fff;
            padding: 15px;
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .player-url {
            font-size: 14px;
            color: #666;
            margin-bottom: 10px;
            word-break: break-all;
            font-family: monospace;
        }
        #errorMessage {
            color: red;
            margin-top: 10px;
            padding: 10px;
            border-radius: 4px;
            background-color: #ffebee;
            display: none;
        }
        .controls {
            margin-top: 10px;
            display: flex;
            gap: 10px;
            justify-content: flex-end;
        }
        .control-button {
            background-color: #4CAF50;
            color: white;
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        .control-button.clear {
            background-color: #f44336;
        }
        .control-button.clear:hover {
            background-color: #d32f2f;
        }
        .player-controls {
            display: flex;
            justify-content: flex-end;
            margin-top: 10px;
            gap: 10px;
        }
        .hidden {
            display: none !important;
        }
        .footer {
            margin-top: 20px;
            text-align: center;
            color: #666;
            font-size: 14px;
            padding: 10px 0;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>FLV 播放器</h1>
        <div class="input-group {{if .HideInput}}hidden{{end}}">
            <textarea id="flvUrls" placeholder="请输入 FLV 视频 URL（每行一个）">{{.URL}}</textarea>
            <div class="controls">
                <button class="control-button clear" onclick="clearPlayers()">清空</button>
                <button onclick="playVideos()">播放全部</button>
            </div>
        </div>
        <div id="errorMessage"></div>
        <div id="players" class="player-container">
        </div>
        <div class="footer">
            系统运维部驱动
        </div>
    </div>

    <script>
        let players = [];

        function showError(message) {
            const errorElement = document.getElementById('errorMessage');
            errorElement.textContent = message;
            errorElement.style.display = 'block';
        }

        function clearPlayers() {
            destroyAllPlayers();
            document.getElementById('flvUrls').value = '';
            document.getElementById('errorMessage').style.display = 'none';
        }

        function createPlayerElement(url, index) {
            const wrapper = document.createElement('div');
            wrapper.className = 'player-wrapper';

            const urlDiv = document.createElement('div');
            urlDiv.className = 'player-url';
            urlDiv.textContent = url;
            wrapper.appendChild(urlDiv);

            const video = document.createElement('video');
            video.id = 'videoElement_' + index;
            video.controls = true;
            video.muted = true;
            wrapper.appendChild(video);

            const controls = document.createElement('div');
            controls.className = 'player-controls';

            const muteButton = document.createElement('button');
            muteButton.className = 'control-button';
            muteButton.textContent = '开启/关闭声音';
            muteButton.onclick = () => toggleMute(video);
            controls.appendChild(muteButton);

            wrapper.appendChild(controls);
            return { wrapper, video };
        }

        function toggleMute(videoElement) {
            videoElement.muted = !videoElement.muted;
        }

        function destroyAllPlayers() {
            players.forEach(player => {
                if (player.flvPlayer) {
                    player.flvPlayer.destroy();
                }
            });
            players = [];
            document.getElementById('players').innerHTML = '';
        }

        function playVideos() {
            const urls = document.getElementById('flvUrls').value.trim().split('\n').filter(url => url.trim());

            if (urls.length === 0) {
                showError('请输入至少一个 FLV 视频 URL');
                return;
            }

            if (!flvjs.isSupported()) {
                showError('您的浏览器不支持 FLV.js，请使用现代浏览器');
                return;
            }

            destroyAllPlayers();
            const playersContainer = document.getElementById('players');
            const errorElement = document.getElementById('errorMessage');
            errorElement.style.display = 'none';

            urls.forEach((url, index) => {
                try {
                    const { wrapper, video } = createPlayerElement(url, index);
                    playersContainer.appendChild(wrapper);

                    const flvPlayer = flvjs.createPlayer({
                        type: 'flv',
                        url: url.trim(),
                        cors: true
                    });

                    flvPlayer.attachMediaElement(video);

                    flvPlayer.on(flvjs.Events.ERROR, (errorType, errorDetail) => {
                        console.error('Player ' + index + ' Error:', errorType, errorDetail);
                        showError('播放器 ' + (index + 1) + ' 出错: ' + errorType);
                    });

                    flvPlayer.load();
                    video.play().catch(e => {
                        console.error('Player ' + index + ' play error:', e);
                        showError('播放器 ' + (index + 1) + ' 播放失败: ' + e.message);
                    });

                    players.push({ flvPlayer, video });
                } catch (error) {
                    console.error('Error initializing player ' + index + ':', error);
                    showError('播放器 ' + (index + 1) + ' 初始化失败: ' + error.message);
                }
            });
        }

        // 页面关闭时清理所有播放器
        window.onbeforeunload = function() {
            destroyAllPlayers();
        };

        // 如果URL参数中有视频地址，自动开始播放
        window.onload = function() {
            const urlInput = document.getElementById('flvUrls');
            if (urlInput && urlInput.value.trim()) {
                playVideos();
            }
        };
    </script>
</body>
</html>`

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(CORSMiddleware())

	// 主页路由
	r.GET("/", func(c *gin.Context) {
		t := template.Must(template.New("index").Parse(htmlTemplate))
		c.Header("Content-Type", "text/html")
		t.Execute(c.Writer, gin.H{
			"URL":       "",
			"HideInput": false,
		})
	})

	// 视频播放路由
	r.GET("/video", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.Redirect(302, "/")
			return
		}

		t := template.Must(template.New("video").Parse(htmlTemplate))
		c.Header("Content-Type", "text/html")
		t.Execute(c.Writer, gin.H{
			"URL":       url,
			"HideInput": true,
		})
	})

	// 启动服务器
	r.Run(":8080")
}
