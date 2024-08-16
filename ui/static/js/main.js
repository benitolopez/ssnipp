// Initialize highlight.js
hljs.highlightAll();

// Copy to clipboard
let text = (document.getElementById("snippet").innerText);

const copyContent = async () => {
  try {
    await navigator.clipboard.writeText(text);
    showCopyMessage();
  } catch (err) {
    showCopyMessage("error");
  }
}
  
document.getElementById("copy-button").addEventListener("click", copyContent);

// Show copy to cliboard message
showCopyMessage = (messageType = "success") => {
  const copyMessage = document.createElement("span");

  copyMessage.style.display = "none";
  copyMessage.innerHTML = messageType === "success" ? "Copied to clipboard!" : "Failed to copy to clipboard!";
  copyMessage.id = "copy-message";
  copyMessage.className = messageType === "success" ? "fixed bottom-1 right-1 p-4 bg-green-100 text-green-500 rounded text-sm" : "fixed bottom-1 right-1 p-4 bg-red-100 text-red-500 rounded text-sm";
  
  document.body.appendChild(copyMessage);

  document.getElementById("copy-message").style.display = "block";

  setTimeout(() => {
    document.getElementById("copy-message").style.display = "none";
    document.getElementById("copy-message").remove();
  }, 2000);
}
  
