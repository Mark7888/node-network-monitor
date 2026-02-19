import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import MainContent from './MainContent';
import Footer from './Footer';
import BottomNav from './BottomNav';

/**
 * Main layout wrapper with sidebar and header
 * Responsive: Desktop shows sidebar, mobile shows bottom navigation
 */
export default function Layout() {
  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar - Desktop only */}
      <Sidebar />
      
      {/* Main content area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header />
        <MainContent>
          <Outlet />
        </MainContent>
        <Footer />
      </div>

      {/* Bottom Navigation - Mobile only */}
      <BottomNav />
    </div>
  );
}
