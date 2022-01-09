import React, { Component } from "react";

class Logout extends React.Component {
  componentDidMount() {
    localStorage.removeItem("jwt");
    window.location = "/login";
  }
  render() {
    return null;
  }
}

export default Logout;
