import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import MainContent from './MainContent';
import Footer from './Footer';

/**
 * Main layout wrapper with sidebar and header
 */
export default function Layout() {
  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <Sidebar />
      
      {/* Main content area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header />
        <MainContent>
          <Outlet />
        </MainContent>
        <Footer />
      </div>
    </div>
  );
}
