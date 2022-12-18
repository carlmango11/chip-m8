import { FunctionComponent, useEffect } from "react";

interface Props {
  onPressed: (n: number) => void;
  onUnpressed: (n: number) => void;
}

export const Keyboard: FunctionComponent<Props> = ({
  onPressed,
  onUnpressed,
}) => {
  useEffect(() => {
    const keyPressed = keyHandler(onPressed);
    const keyUnpressed = keyHandler(onUnpressed);

    document.addEventListener("keydown", keyPressed);
    document.addEventListener("keyup", keyUnpressed);

    return () => {
      document.removeEventListener("keydown", keyPressed);
      document.removeEventListener("keyup", keyUnpressed);
    };
  }, []);

  const numbers = [];
  for (let i = 0; i < 4; i++) {
    let row = [];

    for (let j = 0; j < 4; j++) {
      const n = i * 4 + j;
      const label = n.toString(16).toUpperCase();

      row.push(
        <div
          className="key"
          onMouseDown={() => onPressed(n)}
          onMouseUp={() => onUnpressed(n)}
        >
          {label}
        </div>
      );
    }
    numbers.push(<div className="row">{row}</div>);
  }

  return (
    <div className="keyboard">
      <div className="wrapper">{numbers}</div>
    </div>
  );
};

const keyHandler = (callback: (n: number) => void) => {
  return (event: KeyboardEvent) => {
    const n = parseInt(event.key, 16);
    if (n >= 0 && n <= 15) {
      callback(n);
    }
  };
};
