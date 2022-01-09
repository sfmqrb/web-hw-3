import { Link } from "react-router-dom";
import React from "react";

const ButtonAddNote = () => {
  return (
    <div>
      <Link
        to="/notes/new"
        className="btn btn-primary float-right"
        style={{
          marginBottom: 20,
          marginLeft: 20,
        }}>
        New Note
      </Link>
    </div>
  );
};

export default ButtonAddNote;
