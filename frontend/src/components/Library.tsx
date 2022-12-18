import { FunctionComponent, useEffect, useState } from "react";
import { MenuItem, Select, SelectChangeEvent } from "@mui/material";

interface Props {
  onSelect: (script: string) => void;
}

interface LibraryData {
  [name: string]: string;
}

export const Library: FunctionComponent<Props> = ({ onSelect }) => {
  const [library, setLibrary] = useState<LibraryData>({});
  const [romName, setRomName] = useState<string>("");

  const handleChange = (e: SelectChangeEvent) => {
    const name = e.target.value;

    setRomName(name);
    onSelect(library[name]);
  };

  useEffect(() => {
    const fetchLibrary = async () => {
      const resp = await fetch("/library");
      setLibrary(await resp.json());
    };

    fetchLibrary();
  }, []);

  let items = [];
  for (const name in library) {
    items.push(
      <MenuItem key={name} value={name}>
        {name}
      </MenuItem>
    );
  }

  return (
    <Select
      className="library"
      label="Library"
      value={romName}
      onChange={handleChange}
    >
      {items}
    </Select>
  );
};
