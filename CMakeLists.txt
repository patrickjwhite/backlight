cmake_minimum_required(VERSION 2.8.11)

project(backlight)

set(CMAKE_AUTOMOC ON)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(Qt5Quick REQUIRED)
find_package(Qt5Widgets REQUIRED)

file(GLOB Sources *.cpp)
file(GLOB Resources *.qrc)

include_directories(qt-solutions/qtsingleapplication/src)
set(QtSingleApplicationSources
    qt-solutions/qtsingleapplication/src/qtlocalpeer.cpp
    qt-solutions/qtsingleapplication/src/qtsingleapplication.cpp
    qt-solutions/qtsingleapplication/src/qtsinglecoreapplication.cpp)

qt5_add_resources(ResourceSources ${Resources})

add_executable(backlight ${Sources} ${ResourceSources} ${QtSingleApplicationSources})

target_link_libraries(backlight Qt5::Widgets Qt5::Quick)
