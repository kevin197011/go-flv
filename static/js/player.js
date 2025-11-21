let players = [];

// 统计信息类
class PlayerStats {
  constructor() {
    this.reset();
  }

  reset() {
    this.totalBytes = 0;
    this.totalPackets = 0;
    this.videoFrames = 0;
    this.keyFrames = 0;
    this.lastKeyFrameTime = 0;
    this.currentGOP = 0;
    this.maxGOP = 0;
    this.startTime = Date.now();
    this.lastUpdateTime = Date.now();
    this.lastBytes = 0;
    this.lastFrames = 0;
    this.bitrate = 0; // bps
    this.fps = 0;
  }

  addBytes(bytes) {
    this.totalBytes += bytes;
    this.totalPackets++;
  }

  addVideoFrame(isKeyFrame = false) {
    this.videoFrames++;
    if (isKeyFrame) {
      this.keyFrames++;
      if (this.lastKeyFrameTime > 0) {
        this.currentGOP = this.videoFrames - this.lastKeyFrameTime;
        this.maxGOP = Math.max(this.maxGOP, this.currentGOP);
      }
      this.lastKeyFrameTime = this.videoFrames;
      this.currentGOP = 0;
    } else {
      this.currentGOP++;
    }
  }

  update() {
    const now = Date.now();
    const timeDiff = (now - this.lastUpdateTime) / 1000; // 秒
    const bytesDiff = this.totalBytes - this.lastBytes;
    const framesDiff = this.videoFrames - this.lastFrames;

    if (timeDiff > 0) {
      // 计算码率 (bps)
      this.bitrate = Math.round((bytesDiff * 8) / timeDiff);
      // 计算帧率 (fps)
      this.fps = Math.round((framesDiff / timeDiff) * 10) / 10;
    }

    this.lastUpdateTime = now;
    this.lastBytes = this.totalBytes;
    this.lastFrames = this.videoFrames;
  }

  getStats() {
    this.update();
    return {
      bitrate: this.bitrate,
      fps: this.fps,
      packets: this.totalPackets,
      keyFrames: this.keyFrames,
      gop: this.currentGOP,
      maxGOP: this.maxGOP,
      totalBytes: this.totalBytes,
      totalFrames: this.videoFrames
    };
  }
}

function showNotification(message, type = 'error') {
  const container = document.getElementById('notification-container');
  const id = Date.now();

  const notification = document.createElement('div');
  notification.className = `max-w-sm w-full shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden border transform transition-all duration-300 ease-out ${
    type === 'success' ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'
  }`;
  notification.style.transform = 'translateY(-100%) translateX(100%)';
  notification.style.opacity = '0';

  notification.innerHTML = `
    <div class="p-4">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          ${type === 'success' ?
            '<svg class="h-6 w-6 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>' :
            '<svg class="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>'
          }
        </div>
        <div class="ml-3 w-0 flex-1 pt-0.5">
          <p class="text-sm font-medium ${type === 'success' ? 'text-green-900' : 'text-red-900'}">${message}</p>
        </div>
        <div class="ml-4 flex-shrink-0 flex">
          <button onclick="removeNotification(${id})" class="rounded-md inline-flex focus:outline-none focus:ring-2 focus:ring-offset-2 ${
            type === 'success' ? 'text-green-400 hover:text-green-500' : 'text-red-400 hover:text-red-500'
          }">
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>
  `;

  notification.id = `notification-${id}`;
  container.appendChild(notification);

  // Animate in
  requestAnimationFrame(() => {
    notification.style.transform = 'translateY(0) translateX(0)';
    notification.style.opacity = '1';
  });

  // Auto remove after 4 seconds
  setTimeout(() => {
    removeNotification(id);
  }, 4000);
}

function removeNotification(id) {
  const notification = document.getElementById(`notification-${id}`);
  if (notification) {
    notification.style.transform = 'translateY(-100%) translateX(100%)';
    notification.style.opacity = '0';
    setTimeout(() => {
      if (notification.parentNode) {
        notification.parentNode.removeChild(notification);
      }
    }, 300);
  }
}

function showError(message) {
  showNotification(message, 'error');
}

function showSuccess(message) {
  showNotification(message, 'success');
}

function updatePlayerCount() {
  const countElement = document.getElementById("playerCount");
  if (countElement) {
    countElement.textContent = `(${players.length})`;
  }
}

function updateEmptyState() {
  const emptyState = document.getElementById("emptyState");
  const playersContainer = document.getElementById("players");

  if (players.length === 0) {
    emptyState.style.display = "block";
    playersContainer.style.display = "none";
  } else {
    emptyState.style.display = "none";
    playersContainer.style.display = "grid";
  }
}

function clearPlayers() {
  destroyAllPlayers();
  document.getElementById("flvUrls").value = "";
  showNotification("播放器已清空", "success");
}

function createPlayerElement(url, index) {
  const wrapper = document.createElement("div");
  wrapper.className = "bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-xl transition-shadow duration-200";

  // Header with URL and index
  const header = document.createElement("div");
  header.className = "bg-gradient-to-r from-blue-500 to-blue-600 px-4 py-3";

  const headerContent = document.createElement("div");
  headerContent.className = "flex items-center justify-between";

  const title = document.createElement("h4");
  title.className = "text-white font-semibold flex items-center";
  title.innerHTML = `
    <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
    </svg>
    播放器 ${index + 1}
  `;

  const badge = document.createElement("span");
  badge.className = "px-2 py-1 bg-white bg-opacity-20 rounded-full text-xs text-white";
  badge.textContent = "FLV";

  headerContent.appendChild(title);
  headerContent.appendChild(badge);
  header.appendChild(headerContent);
  wrapper.appendChild(header);

  // URL display
  const urlDiv = document.createElement("div");
  urlDiv.className = "px-4 py-2 text-xs text-gray-600 font-mono break-all bg-gray-50 border-b border-gray-200";
  urlDiv.innerHTML = `
    <div class="flex items-center">
      <svg class="w-4 h-4 mr-2 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
      </svg>
      <span class="truncate">${url}</span>
    </div>
  `;
  wrapper.appendChild(urlDiv);

  // Video wrapper
  const videoWrapper = document.createElement("div");
  videoWrapper.className = "aspect-video bg-black relative";

  const video = document.createElement("video");
  video.id = "videoElement_" + index;
  video.className = "w-full h-full object-contain";
  video.controls = true;
  video.muted = true;

  // Loading indicator
  const loadingDiv = document.createElement("div");
  loadingDiv.className = "absolute inset-0 flex items-center justify-center bg-black bg-opacity-50";
  loadingDiv.innerHTML = `
    <div class="text-white text-center">
      <svg class="animate-spin h-8 w-8 mx-auto mb-2" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="text-sm">加载中...</p>
    </div>
  `;

  videoWrapper.appendChild(video);
  videoWrapper.appendChild(loadingDiv);
  wrapper.appendChild(videoWrapper);

  // Statistics display
  const statsDiv = document.createElement("div");
  statsDiv.id = `stats_${index}`;
  statsDiv.className = "px-4 py-2 bg-gray-900 text-white text-xs font-mono border-t border-gray-700";
  statsDiv.style.display = "none";
  statsDiv.innerHTML = `
    <div class="grid grid-cols-3 gap-2">
      <div><span class="text-gray-400">码率:</span> <span class="text-green-400" id="bitrate_${index}">0</span> bps</div>
      <div><span class="text-gray-400">帧率:</span> <span class="text-green-400" id="fps_${index}">0</span> fps</div>
      <div><span class="text-gray-400">数据包:</span> <span class="text-blue-400" id="packets_${index}">0</span></div>
      <div><span class="text-gray-400">关键帧:</span> <span class="text-yellow-400" id="keyframes_${index}">0</span></div>
      <div><span class="text-gray-400">GOP:</span> <span class="text-purple-400" id="gop_${index}">0</span></div>
      <div><span class="text-gray-400">最大GOP:</span> <span class="text-purple-400" id="maxgop_${index}">0</span></div>
    </div>
  `;
  wrapper.appendChild(statsDiv);

  // Controls
  const controls = document.createElement("div");
  controls.className = "px-4 py-3 bg-gray-50 flex justify-between items-center border-t border-gray-200";

  const statusDiv = document.createElement("div");
  statusDiv.className = "flex items-center text-sm text-gray-600";
  statusDiv.innerHTML = `
    <span class="w-2 h-2 bg-gray-400 rounded-full mr-2"></span>
    <span>准备就绪</span>
  `;

  const muteButton = document.createElement("button");
  muteButton.className = "inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-lg text-gray-700 bg-white hover:bg-gray-50 hover:text-primary transition-colors duration-200";
  muteButton.innerHTML = `
    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 14.142M6.343 6.343L4.929 4.93A1 1 0 003.515 6.343l1.414 1.414L6.343 6.343zM12 2l3 3-3 3-3-3 3-3z"></path>
    </svg>
    声音
  `;
  muteButton.onclick = () => toggleMute(video, muteButton, statusDiv);

  controls.appendChild(statusDiv);
  controls.appendChild(muteButton);
  wrapper.appendChild(controls);

  // Hide loading when video can play
  video.addEventListener('canplay', () => {
    loadingDiv.style.display = 'none';
    statusDiv.innerHTML = `
      <span class="w-2 h-2 bg-green-400 rounded-full mr-2"></span>
      <span>播放中</span>
    `;
    // 显示统计信息
    statsDiv.style.display = 'block';
  });

  video.addEventListener('error', () => {
    loadingDiv.innerHTML = `
      <div class="text-white text-center">
        <svg class="h-8 w-8 mx-auto mb-2 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <p class="text-sm">加载失败</p>
      </div>
    `;
    statusDiv.innerHTML = `
      <span class="w-2 h-2 bg-red-400 rounded-full mr-2"></span>
      <span>播放错误</span>
    `;
  });

  return { wrapper, video, statsDiv };
}

function toggleMute(videoElement, button, statusDiv) {
  videoElement.muted = !videoElement.muted;

  if (videoElement.muted) {
    button.innerHTML = `
      <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" clip-rule="evenodd"></path>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"></path>
      </svg>
      静音
    `;
  } else {
    button.innerHTML = `
      <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 14.142M6.343 6.343L4.929 4.93A1 1 0 003.515 6.343l1.414 1.414L6.343 6.343zM12 2l3 3-3 3-3-3 3-3z"></path>
      </svg>
      声音
    `;
  }
}

function updateStats(playerIndex, stats) {
  const statsData = stats.getStats();
  const formatBitrate = (bps) => {
    if (bps >= 1000000) return (bps / 1000000).toFixed(2) + ' Mbps';
    if (bps >= 1000) return (bps / 1000).toFixed(2) + ' Kbps';
    return bps + ' bps';
  };

  const bitrateEl = document.getElementById(`bitrate_${playerIndex}`);
  const fpsEl = document.getElementById(`fps_${playerIndex}`);
  const packetsEl = document.getElementById(`packets_${playerIndex}`);
  const keyframesEl = document.getElementById(`keyframes_${playerIndex}`);
  const gopEl = document.getElementById(`gop_${playerIndex}`);
  const maxgopEl = document.getElementById(`maxgop_${playerIndex}`);

  if (bitrateEl) bitrateEl.textContent = formatBitrate(statsData.bitrate);
  if (fpsEl) fpsEl.textContent = statsData.fps.toFixed(1);
  if (packetsEl) packetsEl.textContent = statsData.packets;
  if (keyframesEl) keyframesEl.textContent = statsData.keyFrames;
  if (gopEl) gopEl.textContent = statsData.gop;
  if (maxgopEl) maxgopEl.textContent = statsData.maxGOP;
}

function destroyAllPlayers() {
  players.forEach((player) => {
    if (player.flvPlayer) {
      player.flvPlayer.destroy();
    }
    if (player.statsInterval) {
      clearInterval(player.statsInterval);
    }
    if (player.statsCheckInterval) {
      clearInterval(player.statsCheckInterval);
    }
    if (player.performanceCheckInterval) {
      clearInterval(player.performanceCheckInterval);
    }
    if (player.performanceObserver) {
      try {
        player.performanceObserver.disconnect();
      } catch (e) {
        // 忽略错误
      }
    }
  });
  players = [];
  document.getElementById("players").innerHTML = "";
  updatePlayerCount();
  updateEmptyState();
}

function playVideos() {
  const urls = document
    .getElementById("flvUrls")
    .value.trim()
    .split("\n")
    .map((url) => url.trim())
    .filter((url) => url && !url.startsWith("#")); // 过滤空行和以 # 开头的注释行

  if (urls.length === 0) {
    showError("请输入至少一个 FLV 视频 URL");
    return;
  }

  if (!flvjs.isSupported()) {
    showError("您的浏览器不支持 FLV.js，请使用现代浏览器");
    return;
  }

  destroyAllPlayers();
  const playersContainer = document.getElementById("players");

  urls.forEach((url, index) => {
    try {
      const { wrapper, video, statsDiv } = createPlayerElement(url, index);
      playersContainer.appendChild(wrapper);

      // 创建统计对象
      const stats = new PlayerStats();

      const flvPlayer = flvjs.createPlayer({
        type: "flv",
        url: url.trim(),
        cors: true,
      });

      flvPlayer.attachMediaElement(video);

      // 监听统计信息
      let lastVideoTime = 0;
      let lastFrameUpdateTime = Date.now();
      let frameUpdateCount = 0;
      let frameCounter = 0;
      const estimatedKeyFrameInterval = 30; // 假设每30帧有一个关键帧（可根据实际情况调整）
      let lastPlayTime = video.currentTime || 0;

      // 监听视频时间更新来计算帧率和检测关键帧
      video.addEventListener('timeupdate', () => {
        const currentTime = video.currentTime;
        const now = Date.now();
        if (currentTime !== lastVideoTime) {
          frameUpdateCount++;
          frameCounter++;
          lastVideoTime = currentTime;

          // 估算关键帧（基于帧计数）
          if (frameCounter >= estimatedKeyFrameInterval) {
            stats.addVideoFrame(true); // 关键帧
            frameCounter = 0;
          } else {
            stats.addVideoFrame(false); // 非关键帧
          }

          // 每1秒计算一次帧率
          if (now - lastFrameUpdateTime >= 1000) {
            const timeDiff = (now - lastFrameUpdateTime) / 1000;
            const estimatedFps = frameUpdateCount / timeDiff;
            frameUpdateCount = 0;
            lastFrameUpdateTime = now;
          }
        }
      });

      // 监听视频播放来检测关键帧
      video.addEventListener('play', () => {
        lastPlayTime = video.currentTime;
      });

      video.addEventListener('seeked', () => {
        // 时间跳跃通常发生在关键帧
        const timeDiff = Math.abs(video.currentTime - lastPlayTime);
        if (timeDiff > 0.1) {
          stats.addVideoFrame(true); // 检测到关键帧
          frameCounter = 0; // 重置帧计数器
        }
        lastPlayTime = video.currentTime;
      });

      // 监听 flv.js 的统计信息（如果可用）
      const statsCheckInterval = setInterval(() => {
        try {
          // 尝试从 flv.js 获取统计信息
          if (flvPlayer.statisticsInfo) {
            const info = flvPlayer.statisticsInfo;
            if (info && info.totalBytesLoaded) {
              const bytesDiff = info.totalBytesLoaded - stats.totalBytes;
              if (bytesDiff > 0) {
                stats.addBytes(bytesDiff);
              }
            }
          }
        } catch (e) {
          // 忽略错误
        }
      }, 1000);

      // 定期更新统计信息显示
      const statsUpdateInterval = setInterval(() => {
        updateStats(index, stats);
      }, 500); // 每500ms更新一次

      // 监听关键帧（通过 flv.js 事件，如果可用）
      flvPlayer.on(flvjs.Events.MEDIA_INFO, (mediaInfo) => {
        // 媒体信息更新
        if (mediaInfo && mediaInfo.hasVideo) {
          // 可以在这里获取视频信息
        }
      });

      flvPlayer.on(flvjs.Events.ERROR, (errorType, errorDetail) => {
        console.error("Player " + index + " Error:", errorType, errorDetail);
        showError("播放器 " + (index + 1) + " 出错: " + errorType);
      });

      flvPlayer.load();
      video.play().catch((e) => {
        console.error("Player " + index + " play error:", e);
        showError("播放器 " + (index + 1) + " 播放失败: " + e.message);
      });

      // 使用 Performance API 来监控网络传输
      let performanceObserver = null;
      try {
        performanceObserver = new PerformanceObserver((list) => {
          for (const entry of list.getEntries()) {
            if (entry.name && (entry.name.includes(url.trim()) || entry.initiatorType === 'video' || entry.initiatorType === 'xmlhttprequest')) {
              if (entry.transferSize && entry.transferSize > 0) {
                stats.addBytes(entry.transferSize);
              }
            }
          }
        });
        performanceObserver.observe({ entryTypes: ['resource'] });
      } catch (e) {
        // 浏览器可能不支持 Performance Observer API
        console.log('Performance Observer not supported');
      }

      // 定期检查 Performance API 中的资源条目
      let lastCheckedResources = new Set();
      const performanceCheckInterval = setInterval(() => {
        try {
          const resources = performance.getEntriesByType('resource');
          for (const resource of resources) {
            if (resource.name && resource.name.includes(url.trim())) {
              const resourceId = resource.name + '_' + resource.startTime;
              if (!lastCheckedResources.has(resourceId)) {
                lastCheckedResources.add(resourceId);
                if (resource.transferSize && resource.transferSize > 0) {
                  stats.addBytes(resource.transferSize);
                }
              }
            }
          }
        } catch (e) {
          // 忽略错误
        }
      }, 2000);

      players.push({
        flvPlayer,
        video,
        stats,
        statsInterval: statsUpdateInterval,
        statsCheckInterval: statsCheckInterval,
        performanceCheckInterval: performanceCheckInterval,
        performanceObserver
      });
    } catch (e) {
      console.error("Player " + index + " creation error:", e);
      showError("播放器 " + (index + 1) + " 创建失败: " + e.message);
    }
  });

  updatePlayerCount();
  updateEmptyState();
}

// 切换注释状态（注释/取消注释）
function toggleComment() {
  const textarea = document.getElementById("flvUrls");
  if (!textarea) return;

  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const text = textarea.value;
  const lines = text.split("\n");

  // 确定需要处理的行范围
  let startLine = 0;
  let endLine = lines.length - 1;
  let currentPos = 0;

  // 找到选中文本所在的行
  for (let i = 0; i < lines.length; i++) {
    const lineLength = lines[i].length + 1; // +1 for newline
    if (currentPos <= start && start < currentPos + lineLength) {
      startLine = i;
    }
    if (currentPos <= end && end <= currentPos + lineLength) {
      endLine = i;
      break;
    }
    currentPos += lineLength;
  }

  // 如果没有选中文本，只处理当前行
  if (start === end) {
    endLine = startLine;
  }

  // 检查选中行的注释状态
  let allCommented = true;
  let allUncommented = true;
  for (let i = startLine; i <= endLine; i++) {
    const trimmedLine = lines[i].trim();
    if (trimmedLine && !trimmedLine.startsWith("#")) {
      allCommented = false;
    } else if (trimmedLine && trimmedLine.startsWith("#")) {
      allUncommented = false;
    }
  }

  // 决定是注释还是取消注释
  const shouldComment = !allCommented;

  // 记录每行的偏移量变化
  const lineOffsets = new Array(lines.length).fill(0);

  // 处理每一行
  for (let i = startLine; i <= endLine; i++) {
    const line = lines[i];
    const trimmedLine = line.trim();

    if (!trimmedLine) {
      // 空行跳过
      continue;
    }

    if (shouldComment) {
      // 添加注释
      if (!trimmedLine.startsWith("#")) {
        // 找到第一个非空白字符的位置
        const firstNonSpace = line.search(/\S/);
        if (firstNonSpace === -1) {
          lines[i] = "# " + line;
          lineOffsets[i] = 2; // # 和空格
        } else {
          lines[i] = line.slice(0, firstNonSpace) + "# " + line.slice(firstNonSpace);
          lineOffsets[i] = 2; // # 和空格
        }
      }
    } else {
      // 取消注释
      if (trimmedLine.startsWith("#")) {
        const firstNonSpace = line.search(/\S/);
        if (firstNonSpace !== -1 && line[firstNonSpace] === "#") {
          // 移除 # 和后面的空格（如果有）
          let removeLength = 1;
          if (line[firstNonSpace + 1] === " ") {
            removeLength = 2;
          }
          lines[i] = line.slice(0, firstNonSpace) + line.slice(firstNonSpace + removeLength);
          lineOffsets[i] = -removeLength;
        }
      }
    }
  }

  // 计算光标位置的总偏移量
  // 累加从开始到选中位置的所有行的偏移
  let startOffset = 0;
  let endOffset = 0;
  let pos = 0;

  for (let i = 0; i < lines.length; i++) {
    const originalLineLength = i < lines.length - 1
      ? (lines[i].length - lineOffsets[i]) + 1  // 原始长度 + 换行符
      : (lines[i].length - lineOffsets[i]);     // 最后一行没有换行符

    // 计算到 start 位置的偏移
    if (pos + originalLineLength <= start) {
      // 光标在这行之后，累加这行的偏移
      startOffset += lineOffsets[i];
    } else if (pos <= start) {
      // 光标在这行内，累加这行的偏移
      startOffset += lineOffsets[i];
    }

    // 计算到 end 位置的偏移
    if (pos + originalLineLength <= end) {
      endOffset += lineOffsets[i];
    } else if (pos <= end) {
      endOffset += lineOffsets[i];
    }

    pos += originalLineLength;
  }

  // 更新文本内容
  const newText = lines.join("\n");
  textarea.value = newText;

  // 恢复光标位置
  const newStart = Math.max(0, start + startOffset);
  const newEnd = Math.max(0, end + endOffset);
  textarea.setSelectionRange(newStart, newEnd);
  textarea.focus();
}

// 页面关闭时清理所有播放器
window.onbeforeunload = function () {
  destroyAllPlayers();
};

// 如果URL参数中有视频地址，自动开始播放
window.onload = function () {
  updateEmptyState(); // 初始化空状态显示
  const urlInput = document.getElementById("flvUrls");
  if (urlInput && urlInput.value.trim()) {
    playVideos();
  }

  // 添加快捷键支持：Ctrl+/ 或 Cmd+/ 切换注释
  if (urlInput) {
    urlInput.addEventListener("keydown", function (e) {
      // 检测 Ctrl+/ (Windows/Linux) 或 Cmd+/ (Mac)
      if ((e.ctrlKey || e.metaKey) && e.key === "/") {
        e.preventDefault();
        toggleComment();
      }
    });
  }
};
