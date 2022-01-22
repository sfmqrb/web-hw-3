//mohammad
import http from "./httpService";
import { apiUrl } from "../config.json";

const apiEndpoint = apiUrl + "/notes";

function getHeader() {
  return {
    headers: {
      jwt: localStorage.getItem("jwt") || "",

      // jwt: "X",
      "Content-Type": "application/json",
    },
  };
}

function noteUrl(id) {
  return `${apiEndpoint}/${id}`;
}

// export function getNotes() {
//   if (!localStorage.getItem("jwt")) {
//     console.log("first login or register");
//     window.location = "/login";
//     return [];
//   }
//   const localNotes = localStorage.getItem("notes");
//   if (localNotes === null) return JSON.stringify([]);
//   return localNotes;
// }

export function getNotes() {
  if (localStorage.getItem("jwt")) {
    return http.get(apiEndpoint + "/all", getHeader());
  }
  return { data: [] };
  // backend
  // return localStorage.getItem('notes');
}

export function getNote(noteId) {
  console.log("in getNote");
  console.log(getHeader());

  return http.get(noteUrl(noteId), getHeader());
}

export function saveNote(note) {
  console.log("in saveNote");
  if (note._id) {
    const body = { ...note };
    delete body._id;
    return http.put(noteUrl(note._id), body, getHeader());
  }
  return http.post(apiEndpoint + "/new", note, getHeader());
}

export function deleteNote(id) {
  console.log("in deleteNote");
  const notes = JSON.parse(localStorage.getItem("notes")) || [];
  let noteInDb = notes.find((m) => m._id === id);
  notes.splice(notes.indexOf(noteInDb), 1);
  // backend
  return http.delete(noteUrl(id), getHeader());
}
