function createCopyButton(highlightDiv) {
    const button = document.createElement("button");
    button.className = "copy-code-button";
    button.type = "button";
    button.innerText = "Copy";
    button.addEventListener("click", () => {
        let value = highlightDiv.querySelector(":last-child > .chroma > code").innerText
        // replace all the double new lines with single new lines
        value = value.replaceAll("\n\n", "\n");
        navigator.clipboard.writeText(value)
        button.blur();
        button.innerText = "Copied!";
        setTimeout(function () {
            button.innerText = "Copy";
        }, 1000);
    });
    highlightDiv.insertBefore(button, highlightDiv.firstChild);

}

document.querySelectorAll(".highlight").forEach((highlightDiv) => createCopyButton(highlightDiv));
