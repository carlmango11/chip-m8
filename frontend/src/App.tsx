import React, { FunctionComponent, useEffect, useState } from "react";
import "./App.css";
import { Chip8 } from "./components/Chip8";
import { Button, TextField } from "@mui/material";
import "./wasm_exec.js";

import "./App.less";

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
  const [display, setDisplay] = useState<string[]>([]);

  const refreshDisplay = () => {
    setDisplay(window.getDisplay());
  };

  const startChip8 = () => {
    window.load(script.replace(/\s+/g, ""));
  };

  useEffect(() => {
    initWasm();

    const id = setInterval(refreshDisplay, 1000 / FPS);

    return () => {
      clearInterval(id);
    };
  }, []);

  return (
    <div className="app">
      <div className="main-panel">
        <h1>CHIP-8 Virtual Machine</h1>
        <p>Input your program script and click load to start</p>

        <TextField
          className="script-input"
          label="Program Instructions"
          multiline
          rows={20}
          size="small"
          value={script}
          onChange={(v) => setScript(v.target.value)}
        />

        <Button
          className="load-button"
          onClick={startChip8}
          variant="contained"
        >
          Load
        </Button>
      </div>

      <Chip8
        display={display}
        onPressed={window.keyPressed}
        onUnpressed={window.keyUnpressed}
      />
    </div>
  );
};

export default App;
