import React, { Component } from "react";
import NotesTable from "./notesTable";
import ListGroup from "./common/listGroup";
import Pagination from "./common/pagination";
import { deleteNote, getNotes } from "../services/fakeNoteService";
import { getTypes } from "../services/typeService";
import { paginate } from "../utils/paginate";
import ButtonAddNote from "./buttonAddNote";

import _ from "lodash";
import SearchBox from "./searchBox";
import "bootstrap/dist/css/bootstrap-grid.css";
// import { object } from "prop-types";

class Notes extends Component {
  state = {
    notes: [],
    types: [],
    currentPage: 1,
    pageSize: 4,
    searchQuery: "",
    selectedType: null,
    sortColumn: { path: "title", order: "asc" },
  };

  componentDidMount() {
    const types = [{ _id: "", name: "All Types" }, ...getTypes()];
    const x = getNotes();
    // console.log(x[0])
    // console.log(JSON.parse(getNotes()))
    this.setState({ notes: JSON.parse(x), types });
    // backend use cached data to set state notes
    // modify getNotes() to return cached data
    // delete cached notes
    // console.log("delete notes from localStorage");
    // localStorage.removeItem("notes");
  }

  // backend only if getType() is working with backend
  // async componentDidMount() {
  //   const { data } = await getTypes();
  //   const types = [{ _id: "", name: "All Types" }, ...data];
  //   const { data: notes } = await getNotes();
  //   this.setState({ notes, types });
  // }

  handleDelete = async (note) => {
    // locally delete
    const notes = this.state.notes.filter((m) => m._id !== note._id);
    this.setState({ notes });
    // server delete
    console.log(note);
    let result = await deleteNote(note._id);
    console.log(result);
  };

  // handleDelete = async (note) => {
  //   // locally delete
  //   const originalNotes = this.state.notes;
  //   const notes = originalNotes.filter((m) => m._id !== note._id);
  //   this.setState({ notes });
  //   // server delete
  //   try {
  //     await deleteNote(note._id);
  //   } catch (ex) {
  //     if (ex.response && ex.response.status === 404)
  //       toast.error("This note has already been deleted.");
  //     this.setState({ notes: originalNotes });
  //   }
  // };

  handlePageChange = (page) => {
    this.setState({ currentPage: page });
  };

  handleTypeSelect = (type) => {
    this.setState({ selectedType: type, searchQuery: "", currentPage: 1 });
  };

  handleSearch = (query) => {
    this.setState({ searchQuery: query, selectedType: null, currentPage: 1 });
  };

  handleSort = (sortColumn) => {
    this.setState({ sortColumn });
  };

  getPagedData = () => {
    const {
      pageSize,
      currentPage,
      sortColumn,
      selectedType,
      searchQuery,
      notes: allNotes,
    } = this.state;

    let filtered = allNotes;
    if (searchQuery)
      filtered = allNotes.filter(
        (m) =>
          m.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          m.text.toLowerCase().includes(searchQuery.toLowerCase())
      );
    else if (selectedType && selectedType._id)
      filtered = allNotes.filter((m) => m.type._id === selectedType._id);

    const sorted = _.orderBy(filtered, [sortColumn.path], [sortColumn.order]);

    const notes = paginate(sorted, currentPage, pageSize);

    return { totalCount: filtered.length, data: notes };
  };
  render() {
    this.state.notes = this.state.notes === null ? [] : this.state.notes;
    const { length: count } = this.state.notes;
    const { pageSize, currentPage, sortColumn, searchQuery } = this.state;

    if (count === 0)
      return (
        <div>
          <p className="float-left">There are no notes in the database.</p>
          <ButtonAddNote className="float-right" />
        </div>
      );

    const { totalCount, data: notes } = this.getPagedData();

    return (
      <div className="row">
        <div className="col-3 mt-5">
          <ListGroup
            items={this.state.types}
            selectedItem={this.state.selectedType}
            onItemSelect={this.handleTypeSelect}
            style={{ justifyContent: "center" }}
          />
        </div>
        <div className="col">
          <div className="mt-5">
            <ButtonAddNote />
            <p>Showing {totalCount} notes in the database.</p>
          </div>
          <SearchBox value={searchQuery} onChange={this.handleSearch} />
          <NotesTable
            notes={notes}
            sortColumn={sortColumn}
            onLike={this.handleLike}
            onDelete={this.handleDelete}
            onSort={this.handleSort}
          />
          <Pagination
            itemsCount={totalCount}
            pageSize={pageSize}
            currentPage={currentPage}
            onPageChange={this.handlePageChange}
          />
        </div>
      </div>
    );
  }
}

export default Notes;
