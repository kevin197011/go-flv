let players = [];

function showError(message) {
  const errorElement = document.getElementById("errorMessage");
  errorElement.querySelector("h3").textContent = message;
  errorElement.classList.remove("hidden");
}

function clearPlayers() {
  destroyAllPlayers();
  document.getElementById("flvUrls").value = "";
  document.getElementById("errorMessage").classList.add("hidden");
}

function createPlayerElement(url, index) {
  const wrapper = document.createElement("div");
  wrapper.className = "bg-white rounded-lg shadow overflow-hidden";

  const urlDiv = document.createElement("div");
  urlDiv.className = "px-4 py-2 text-sm text-gray-600 font-mono break-all border-b border-gray-200";
  urlDiv.textContent = url;
  wrapper.appendChild(urlDiv);

  const videoWrapper = document.createElement("div");
  videoWrapper.className = "aspect-video bg-black";

  const video = document.createElement("video");
  video.id = "videoElement_" + index;
  video.className = "w-full h-full object-contain";
  video.controls = true;
  video.muted = true;
  videoWrapper.appendChild(video);
  wrapper.appendChild(videoWrapper);

  const controls = document.createElement("div");
  controls.className = "px-4 py-3 bg-gray-50 flex justify-end";

  const muteButton = document.createElement("button");
  muteButton.className = "rounded-md bg-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-dark focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary";
  muteButton.textContent = "开启/关闭声音";
  muteButton.onclick = () => toggleMute(video);
  controls.appendChild(muteButton);

  wrapper.appendChild(controls);
  return { wrapper, video };
}

function toggleMute(videoElement) {
  videoElement.muted = !videoElement.muted;
}

function destroyAllPlayers() {
  players.forEach((player) => {
    if (player.flvPlayer) {
      player.flvPlayer.destroy();
    }
  });
  players = [];
  document.getElementById("players").innerHTML = "";
}

function playVideos() {
  const urls = document
    .getElementById("flvUrls")
    .value.trim()
    .split("\n")
    .filter((url) => url.trim());

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
  const errorElement = document.getElementById("errorMessage");
  errorElement.classList.add("hidden");

  urls.forEach((url, index) => {
    try {
      const { wrapper, video } = createPlayerElement(url, index);
      playersContainer.appendChild(wrapper);

      const flvPlayer = flvjs.createPlayer({
        type: "flv",
        url: url.trim(),
        cors: true,
      });

      flvPlayer.attachMediaElement(video);

      flvPlayer.on(flvjs.Events.ERROR, (errorType, errorDetail) => {
        console.error("Player " + index + " Error:", errorType, errorDetail);
        showError("播放器 " + (index + 1) + " 出错: " + errorType);
      });

      flvPlayer.load();
      video.play().catch((e) => {
        console.error("Player " + index + " play error:", e);
        showError("播放器 " + (index + 1) + " 播放失败: " + e.message);
      });

      players.push({ flvPlayer, video });
    } catch (e) {
      console.error("Player " + index + " creation error:", e);
      showError("播放器 " + (index + 1) + " 创建失败: " + e.message);
    }
  });
}

// 页面关闭时清理所有播放器
window.onbeforeunload = function () {
  destroyAllPlayers();
};

// 如果URL参数中有视频地址，自动开始播放
window.onload = function () {
  const urlInput = document.getElementById("flvUrls");
  if (urlInput && urlInput.value.trim()) {
    playVideos();
  }
};
