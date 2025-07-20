package jsfunctions

import (
	"bytes"
	"fmt"
	"goscreener/internal/model"
)

func AttachTitle() string {
	return `
var titleText = document.querySelector('title').text;
var titleTab = document.querySelector('.fake-nav-tab-title');
if (titleText && titleText != "" && titleTab) {
	titleTab.innerText += titleText;
}
var urlPanel = document.querySelector('.fake-nav-url-panel');
if (urlPanel) {
	urlPanel.innerText += location.protocol + '//' + location.host + location.pathname;
}
`
}

func AttachIcon() string {
	return `
function checkIfImageExists(url) {
  const img = new Image();
  img.src = url;
  
  if (img.complete) {
    return true;
  } else {
    img.onload = () => {
      return true;
    };
    
    img.onerror = () => {
      return false;
    };
  }
}

var icon = location.protocol + '//' + location.host + '/favicon.ico';
if (checkIfImageExists(icon)) {
	document.querySelector('.fake-nav-tab-logo').setAttribute("src", icon);
} else {
	var iconLink = document.querySelector('[rel="icon"]');
	if (iconLink) {
		var ilink = iconLink.getAttribute('href');
		if (ilink && ilink != "") {
		document.querySelector('.fake-nav-tab-logo').setAttribute("src", ilink);
		} 
	}
}
`
}

func MakeFullScreenActions() string {
	return ` function scrollAndCallback(callback) {
    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });

    setTimeout(() => {
        callback();
    }, 1000);
}

scrollAndCallback(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
});
 `
}

func MakeFixedNodes(nodeSelectors []*model.NodeSelector) string {
	if nodeSelectors == nil || len(nodeSelectors) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("(function() {")

	for i, node := range nodeSelectors {
		if node.Many == true {
			buffer.WriteString(fmt.Sprintf(
				` var nodes%d = document.querySelectorAll('%s');
if (nodes%d) {
nodes%d.forEach(function(node) {
node.style.position = 'static';
});
} `,
				i, node.Selector, i, i))
		} else {
			if node.Parent {
				buffer.WriteString(fmt.Sprintf(
					` var node%d = document.querySelector('%s');
if (node%d && node%d.parentElement) {
node%d.parentElement.style.position = 'static';
} `,
					i, node.Selector, i, i, i))
			} else {
				buffer.WriteString(fmt.Sprintf(
					` var node%d = document.querySelector('%s');
if (node%d) {
node%d.style.position = 'static';
} `,
					i, node.Selector, i, i))
			}
		}
	}

	buffer.WriteString("})();")
	return buffer.String()
}

func MakeRemoveNodesScript(nodeSelectors []*model.NodeSelector) string {
	if nodeSelectors == nil || len(nodeSelectors) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("(function() {")

	for i, node := range nodeSelectors {
		if node.Many == true {
			buffer.WriteString(fmt.Sprintf(
				` var nodes%d = document.querySelectorAll('%s');
if (nodes%d) {
nodes%d.forEach(function(node) {
node.remove();
});
} `,
				i, node.Selector, i, i))
		} else {
			if node.Parent {
				buffer.WriteString(fmt.Sprintf(
					` var node%d = document.querySelector('%s');
if (node%d && node%d.parentElement) {
node%d.parentElement.remove();
} `,
					i, node.Selector, i, i, i))
			} else {
				buffer.WriteString(fmt.Sprintf(
					` var node%d = document.querySelector('%s');
if (node%d) {
node%d.remove();
} `,
					i, node.Selector, i, i))
			}
		}
	}

	buffer.WriteString("})();")
	return buffer.String()
}
