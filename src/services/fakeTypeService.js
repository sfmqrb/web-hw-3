export const types = [
  { _id: "5b21ca3eeb7f6fbccd471818", name: "Personal" },
  { _id: "5b21ca3eeb7f6fbccd471814", name: "Work" },
  { _id: "5b21ca3eeb7f6fbccd471820", name: "Family" },
  { _id: "5b21ca3eeb7f6fbccd471821", name: "Others" },
];

export function getTypes() {
  return types.filter((g) => g);
}
