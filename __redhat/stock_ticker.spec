%define debug_package   %{nil}
%define _build_id_links none
%define _name stock_ticker
%define _prefix /opt
%define _version 0.0.1
%define _rel 1
%define _arch x86_64
%define _binaryname stockticker
%define _binaryname_daemon %{_binaryname}Daemon

Name:       stock_ticker
Version:    %{_version}
Release:    %{_rel}
Summary:    Stock Ticker game adaptation

Group:      Games
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/mainline/stock_ticker.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
Stock Ticker game adaptation

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{name}-%{version}/%{_binaryname} ./client
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{name}-%{version}/%{_binaryname_daemon} ./daemon

%clean
rm -rf $RPM_BUILD_ROOT

%pre

%install
rm -rf %{buildroot}
install -Dpm 0755 %{_builddir}/%{name}-%{version}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}
install -Dpm 0755 %{_builddir}/%{name}-%{version}/%{_binaryname_daemon} %{buildroot}%{_sbindir}/%{_binaryname_daemon}

%post

%preun

%postun

%files
%defattr(0755,root,root,-)
%{_bindir}/%{_binaryname}
%{_sbindir}/%{_binaryname_daemon}


%changelog
