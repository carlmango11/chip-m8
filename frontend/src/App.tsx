import React, { FunctionComponent, useEffect, useState } from "react";
import "./App.css";
import { Chip8 } from "./components/Chip8";
import { Button, TextField } from "@mui/material";
import "./wasm_exec.js";

const FPS = 25;

function initWasm() {
  const go = new window.Go();

  WebAssembly.instantiateStreaming(fetch("chip-8.wasm"), go.importObject).then(
    (result) => {
      go.run(result.instance);
    }
  );
}

const App: FunctionComponent = () => {
  const [script, setScript] = useState<string>("");
  const [display, setDisplay] = useState<number[]>([]);

  const refreshDisplay = () => {
    setDisplay(window.getDisplay());
  };

  useEffect(() => {
    initWasm();

    const id = setInterval(refreshDisplay, 1000 / FPS);

    return () => {
      clearInterval(id);
    };
  }, []);

  return (
    <div className="App">
      <TextField
        label="Program Instructions"
        multiline
        maxRows={4}
        value={script}
        onChange={(v) => setScript(v.target.value)}
      />

      <Button variant="contained">Load</Button>

      <Chip8
        display={display}
        onPressed={window.keyPressed}
        onUnpressed={window.keyUnpressed}
      />
    </div>
  );
};

export default App;
