import React from "react";
import FloatText from "./floatText";
import { uuid } from "uuidv4";

const styles = {
  overlay: {
    position: "fixed",
    top: 0,
    left: 0,
    zIndex: 2,
    width: "100%",
    height: "100%",
    pointerEvents: "none"
  }
};

class OverlayWindow extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      floatTextList: []
    };
  }

  deleteText(id) {
    const floatTextList = this.state.floatTextList.filter(
      floatText => floatText.props.id !== id
    );
    this.setState({
      floatTextList: floatTextList
    });
  }

  render() {
    return <div style={styles.overlay}>{this.state.floatTextList}</div>;
  }

  addText(text) {
    const floatTextList = this.state.floatTextList;
    const generatedUuid = uuid();
    console.log(generatedUuid);
    floatTextList.push(
      <FloatText
        key={generatedUuid}
        id={generatedUuid}
        text={text}
        deleteSelf={id => this.deleteText(id)}
      />
    );
    this.setState({
      floatTextList: floatTextList
    });
  }
}

export default OverlayWindow;
