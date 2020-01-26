import React from "react";
import ReactDOM from "react-dom";
import TextField from "@material-ui/core/TextField";
import Box from "@material-ui/core/Box";
import Switch from "@material-ui/core/Switch";

console.log(`Popup`);

class Popup extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      serverUrl: "ws://127.0.0.1:8080/ws", // TODO must be changeable
      active: false
    };
  }

  handleTextChange = name => event => {
    this.setState({ [name]: event.target.document });
  };

  sendMessageToBackground = () => {
    chrome.runtime.sendMessage({
      type: "start",
      data: { wsUrl: this.state.serverUrl }
    });
  };

  toggleActive = () => {
    const current = this.state.active;
    if (current == true) {
    } else {
      this.sendMessageToBackground();
    }
    this.setState({ active: !current });
  };

  render() {
    return (
      <div>
        <Box width={400}>
          <TextField
            value={this.state.serverUrl}
            label="Server Url"
            onChange={this.handleTextChange("serverUrl")}
          />
          <Switch
            checked={this.state.active}
            onChange={this.toggleActive}
            label="Active"
          />
        </Box>
      </div>
    );
  }
}

ReactDOM.render(<Popup />, document.getElementById("app"));
