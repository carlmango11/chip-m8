import { FunctionComponent, useEffect, useState } from "react";
import { MenuItem, Select } from "@mui/material";

interface Props {
  onSelect: (script: string) => void;
}

interface Library {
  [name: string]: string;
}

export const Library: FunctionComponent<Props> = ({ onSelect }) => {
  const [libary, setLibrary] = useState<Library>({});

  useEffect(() => {
    const fetchLibrary = async () => {
      const resp = await fetch("/library");
      setLibrary(await resp.json());
    };

    fetchLibrary();
  }, []);

  let items = [];
  for (const name in libary) {
    items.push(<MenuItem value={libary[name]}>{name}</MenuItem>);
  }

  return (
    <Select
      label="Library"
      onChange={(e) => onSelect(e.target.value as string)}
    >
      {items}
    </Select>
  );
};
