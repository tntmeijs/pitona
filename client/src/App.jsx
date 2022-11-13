import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { PageLoader } from "./components/PageLoader";
import { PageWrapper } from "./components/PageWrapper";
import { ContextProvider } from "./Context";
import { OverviewPage } from "./pages/OverviewPage";
import { RawData } from "./pages/RawData";
import { Speedometer } from "./pages/Speedometer";
import { Temperature } from "./pages/Temperature";

export function App() {
  return (
    <ContextProvider>
      <BrowserRouter>
        <PageLoader />

        <PageWrapper>
          <Routes>
            <Route index element={<OverviewPage />} />
            <Route path="*" element={<Navigate to="/" />} />
            <Route path="raw" element={<RawData />} />
            <Route path="speedometer" element={<Speedometer />} />
            <Route path="temperature" element={<Temperature />} />
          </Routes>
        </PageWrapper>
      </BrowserRouter>
    </ContextProvider>
  );
}
