import React, { Component } from "react";
import { Link } from "react-router-dom";
import NotesTable from "./notesTable";
import ListGroup from "./common/listGroup";
import Pagination from "./common/pagination";
import { getNotes, deleteNote } from "../services/fakeNoteService";
import { getTypes } from "../services/typeService";
import { paginate } from "../utils/paginate";
import _ from "lodash";
import SearchBox from "./searchBox";

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

    this.setState({ notes: getNotes(), types });
  }

  handleDelete = (note) => {
    const notes = this.state.notes.filter((m) => m._id !== note._id);
    this.setState({ notes });

    deleteNote(note._id);
  };

  handleLike = (note) => {
    const notes = [...this.state.notes];
    const index = notes.indexOf(note);
    notes[index] = { ...notes[index] };
    notes[index].liked = !notes[index].liked;
    this.setState({ notes });
  };

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
    const { length: count } = this.state.notes;
    const { pageSize, currentPage, sortColumn, searchQuery } = this.state;

    if (count === 0) return <p>There are no notes in the database.</p>;

    const { totalCount, data: notes } = this.getPagedData();

    return (
      <div className="row">
        <div className="col-3">
          <ListGroup
            items={this.state.types}
            selectedItem={this.state.selectedType}
            onItemSelect={this.handleTypeSelect}
          />
        </div>
        <div className="col">
          <Link
            to="/notes/new"
            className="btn btn-primary"
            style={{ marginBottom: 20 }}>
            New Note
          </Link>
          <p>Showing {totalCount} notes in the database.</p>
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
