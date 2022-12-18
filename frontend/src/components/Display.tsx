import { FunctionComponent } from "react";

interface Props {
  display: string[];
}

export const Display: FunctionComponent<Props> = ({ display }) => {
  const pixels = display.map((row, i) => {
    const rowPixels = [];

    for (let i = 0; i < 64; i++) {
      const className = row[i] === "1" ? "pixel on" : "pixel off";
      rowPixels.push(<span key={i} className={className} />);
    }

    return (
      <span key={i} className="row">
        {rowPixels}
      </span>
    );
  });

  return <div className="display">{pixels}</div>;
};
