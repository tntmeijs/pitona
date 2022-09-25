import { Navbar } from "./components/Navbar";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { DebuggingPage } from "./pages/DebuggingPage";
import { OverviewPage } from "./pages/OverviewPage";
import { SystemPage } from "./pages/SystemPage";
import { PageWrapper } from "./pages/PageWrapper";

export function App() {
  return (
    <BrowserRouter>
      <Navbar />

      <PageWrapper>
        <Routes>
          <Route index element={<OverviewPage />} />
          <Route path="/debugging" element={<DebuggingPage />} />
          <Route path="/system" element={<SystemPage />} />
        </Routes>
      </PageWrapper>
    </BrowserRouter>
  );
}
