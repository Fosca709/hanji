cmake_minimum_required(VERSION 3.21)

project(Hanji VERSION 0.2.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_AUTORCC ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(Qt6 6.2 REQUIRED COMPONENTS DBus Widgets)

qt_add_executable(hanji
    src/main.cpp
    assets/resources.qrc
)

target_link_libraries(hanji PRIVATE Qt6::DBus Qt6::Widgets)

install(TARGETS hanji RUNTIME DESTINATION bin)
