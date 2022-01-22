import React, { Component } from "react";
import { Link } from "react-router-dom";
import Table from "./common/table";

class NotesTable extends Component {
  columns = [
    {
      path: "title",
      label: "Title",
      content: (note) => <Link to={`/notes/${note._id}`}>{note.title}</Link>,
      sizeClass: "col-3",
    },
    {
      path: "type.name",
      content: (note) => note.type || " ",
      sizeClass: "col-3",
      label: "Type",
    },
    {
      path: "text",
      label: "Text",
      content: (note) => note.text,
      sizeClass: "col-12",
    },
    {
      key: "delete",
      content: (note) => (
        <button
          onClick={() => this.props.onDelete(note)}
          className="btn btn-danger btn-sm">
          Delete
        </button>
      ),
      sizeClass: "",
    },
  ];

  render() {
    const { notes, onSort, sortColumn } = this.props;

    return (
      <Table
        columns={this.columns}
        data={notes}
        sortColumn={sortColumn}
        onSort={onSort}
      />
    );
  }
}

export default NotesTable;
