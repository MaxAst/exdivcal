import React from "react";
import CreatableSelect from "react-select/creatable";

function Index() {
  const handleChange = (newValue: any, actionMeta: any) => {
    console.group("Value Changed");
    console.log(newValue);
    console.log(`action: ${actionMeta.action}`);
    console.groupEnd();
  };
  return (
    <main>
      <h1>ExDivCal</h1>
      <h3>
        A service to download an .ics file (compatible with Outlook and Apple's
        Calendar) containing the ex dividend dates of the companies you select.
      </h3>
      <form>
        <label>Please fill in the relevant stock symbols:</label>
        <CreatableSelect isMulti onChange={handleChange} />
      </form>
      <p>
        <b>Note:</b> They need to be the exact same symbols used by{" "}
        <a href="https://finance.yahoo.com">finance.yahoo.com</a>. Please
        manually check that you are filling in the right symbols by getting them
        from{" "}
        <a href="https://finance.yahoo.com/lookup/equity?s=SomeCompanyName">
          https://finance.yahoo.com/lookup/equity?s=SomeCompanyName
        </a>
      </p>
    </main>
  );
}

export default Index;
