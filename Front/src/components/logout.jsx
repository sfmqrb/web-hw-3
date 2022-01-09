import React, { Component } from "react";

class Logout extends React.Component {
  componentDidMount() {
    // localStorage.removeItem("jwt");
    // localStorage.removeItem("notes");
    localStorage.clear();
    window.location = "/login";
  }
  render() {
    return null;
  }
}

export default Logout;
