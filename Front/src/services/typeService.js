export const types = [
  { _id: "Personal", name: "Personal" },
  { _id: "Work", name: "Work" },
  { _id: "Family", name: "Family" },
  { _id: "Others", name: "Others" },
];

export function getTypes() {
  return types.filter((g) => g);
}
