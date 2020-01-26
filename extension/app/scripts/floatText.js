import React from "react";
import PropTypes from "prop-types";

const propTypes = {
  deleteSelf: PropTypes.func
};

const styles = {
  text: {
    position: "absolute",
    color: "#000"
  }
};

class FloatText extends React.Component {
  constructor(props) {
    super(props);

    const text = !props.text ? "test" : props.text;
    this.state = {
      x: -100,
      y: Math.random() * 500,
      dx: 5 + 10 * Math.random(),
      dy: 0,
      remainingLifetime: 10 * 60,
      text: text,
      fontsize: 60,
      frameRate: 60
    };

    this.mainLoop = this.mainLoop.bind(this);

    setTimeout(this.mainLoop, 1000 / this.state.frameRate);
  }

  mainLoop() {
    this.setState({
      x: this.state.x + this.state.dx,
      y: this.state.y + this.state.dy,
      remainingLifetime: this.state.remainingLifetime - 1
    });

    if (this.isOutdated()) {
      this.props.deleteSelf(this.props.id);
      return;
    }

    setTimeout(this.mainLoop, 1000 / this.state.frameRate);
  }

  isOutdated() {
    return this.state.remainingLifetime <= 0;
  }

  getStyle() {
    let ButtonStyles = {
      fontSize: this.state.fontsize + "px",
      top: this.state.y + "px",
      left: this.state.x + "px"
    };
    return Object.assign(ButtonStyles, styles.text);
  }

  render() {
    return (
      <div className="text" style={this.getStyle()}>
        {this.state.text}
      </div>
    );
  }
}

FloatText.propTypes = propTypes;
export default FloatText;
