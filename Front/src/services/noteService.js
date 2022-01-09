import http from "./httpService";
import { apiUrl } from "../config.json";

const apiEndpoint = apiUrl + "/notes";

function noteUrl(id) {
  return `${apiEndpoint}/${id}`;
}

export function getNotes() {
  return http.get(apiEndpoint);
  // backend
  // return localStorage.getItem('notes');
}

export function getNote(noteId) {
  return http.get(noteUrl(noteId));
}

export function saveNote(note) {
  if (note._id) {
    const body = { ...note };
    delete body._id;
    return http.put(noteUrl(note._id), body);
  }

  return http.post(apiEndpoint, note);
}

export function deleteNote(noteId) {
  return http.delete(noteUrl(noteId));
}
