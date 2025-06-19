# Application Requirements Document (ARD)  
**Version:** 2.0  
**Last Updated:** June 19, 2025

---

## 1. Introduction

### 1.1 Purpose  
This document defines the functional and non-functional requirements for a cross-platform text analysis application with an enhanced user interface. The application must efficiently handle large text files (via chunked or full in-memory loading, as appropriate), support dynamic filtering through an integrated filter management system, and provide a modern, user-friendly GUI with split panels and a top menu bar. The interface will be built using Raylib and allow both keyboard and mouse interactions.

### 1.2 Scope  
- **Target Platforms:** Linux, macOS, and Windows.
- **Core Functionality:**  
  - Efficient file loading and processing (full or chunked mode).
  - Dynamic filtering (including multi-filter support).
  - Advanced GUI with dual panels, pop-up dialogs, and menu-based operations.
- **Graphical User Interface (GUI):**  
  - A split-window layout with a top menu bar.
  - The top panel for rendering text.
  - The bottom panel for managing filters.
- **Performance:**  
  - Minimum latency in file operations.
  - Smooth scrolling and rapid UI updates with asynchronous I/O and caching.
- **Extensibility:**  
  - A modular architecture using a controller/state manager that decouples business logic from UI logic.

---

## 2. System Overview

The application is composed of several clearly defined modules:

- **Data Manager / File I/O Manager:**  
  Handles file loading (entire file or in fixed-size chunks), buffering, and caching.

- **Filter Manager:**  
  Manages a list of filters (each with `enabled` status, `filterText`, and ordering) and applies them to the loaded data.

- **State Manager / Controller:**  
  Maintains the overall application state, including mode selection (filtered vs. unfiltered), file offsets, and data caches. Responsible for regenerating data in response to filter changes.

- **UI Manager (GUI Controller):**  
  Implements all GUI functionality using Raylib. Handles layout management (split panels, pop-ups, menus), rendering of text and filter panels, and user interaction events (clicks, right-clicks, menu selections).

- **Event/Notification Manager:**  
  Provides asynchronous and decoupled communication between the UI Manager, Data Manager, Filter Manager, and State Manager.

- **Input/Interaction Manager:**  
  Processes user inputs (keyboard and mouse) and dispatches them appropriately.

- **Configuration Manager:**  
  Centralizes configurable parameters such as default chunk size, default filters, UI preferences, etc.

- **Logging/Error Manager:**  
  Handles logging events for debugging and operation monitoring.

- **Asynchronous Task/Thread Manager (Optional):**  
  Manages background operations to ensure a responsive UI.

---

## 3. Functional Requirements

### 3.1 File Loading and I/O  
- **Chunked File Reading:**  
  - For files larger than 5 MB:  
    - Load data in fixed-size chunks (e.g., 10 MB).
    - Maintain a sliding window of 3 contiguous chunks.
    - Handle partial lines at boundaries by carrying incomplete line data to the next chunk.
- **Full File Loading:**  
  - For files 5 MB or smaller, load the entire file into memory to allow instantaneous operations.

### 3.2 Filtering Engine  
- **Filter Data Structure:**  
  - Each filter object shall contain:  
    - `enabled` (boolean)  
    - `filterText` (string)  
    - An indicator for ordering (to define priority/display order).  
- **Filter Examples:**  
  - Filter 1: enabled = true, filterText = "ERROR"
  - Filter 2: enabled = true, filterText = "log"
  - Filter 3: enabled = false, filterText = "Login"
- **Filter Combination:**  
  - Support multi-filter mode using logical OR (a row is displayed if it contains any enabled filter text) with potential for future support of AND semantics.
- **Dynamic Filter Management:**  
  - Allow creation, deletion, updating (enabled/disabled state, text value), and reordering via the GUI.
  - When filter settings change, trigger regeneration of the data (either by reprocessing cached chunks or re-evaluating the in-memory text), and update the display.

### 3.3 Scrolling and Navigation  
- **Smooth Scrolling:**  
  - Support rapid scrolling across text using a sliding window mechanism of cached chunks.
  - When scrolling beyond the current window:
    - Load the next/previous chunk dynamically.
    - Preload adjacent chunks to enable quick back-and-forth scrolling.
- **UI Responsiveness:**  
  - Maintain smooth transitions through asynchronous I/O and background tasks.

### 3.4 User Interface (GUI)

#### 3.4.1 General Layout  
- **Split Window Layout:**  
  - **Text Panel (Top Panel):**  
    - Primary display area for rendering the text loaded from the file.
    - Supports right-click interaction on any row.
  - **Filter Panel (Bottom Panel):**  
    - Dedicated area for displaying and managing filters.
    - Shows a list of active filters along with their properties.

- **Top Menu Bar:**  
  - Contains menus: **File**, **Edit**, and **Help.**
  - **File Menu Options:**  
    - **Load File:** Open a text file.
    - **Reload File:** Reload the current file.
    - **Save File:** Save a copy of the current file (with applied filters) under a different name.
    - **Close File:** Close the currently open file.
  - **Edit Menu Options:**  
    - **Create Filter:** Open a dialog to create a new filter.
    - **Disable All Filters:** Set all filters' enabled status to false.
    - **Enable All Filters:** Set all filters' enabled status to true.
    - **Export Filters:** Save current filters (and their states) to an external file for future use or sharing.
    - **Load Filters:** Import filters from an external file and load them into the application (the filters are then displayed in the filter panel).
  - **Help Menu Options:**  
    - Access to user documentation, version information, and contact details.

#### 3.4.2 Filter Panel (Bottom Panel) Interactions  
- **Capabilities:**
  - **Create New Filter:**  
    - Button or menu option within the panel to open a "Create New Filter" pop-up dialog.
  - **Delete Filter:**  
    - Option available (e.g., via a context menu or delete button) next to each filter entry.
  - **Update Filter Properties:**  
    - Modify filter status (enabled/disabled) by toggling a checkbox.
    - Edit the text value of the filter.
    - Reorder filters (e.g., via drag-and-drop or up/down controls).
- **Integration with Text Panel:**  
  - **Right-Click Interaction:**  
    - Allow right-click on any row in the text panel.
    - Display a context menu with the option "Create New Filter" pre-populated with the clicked row’s full text.
    - On selection, open the "Create New Filter" pop-up dialog with the row’s text set as the filter criteria.

#### 3.4.3 Create New Filter Pop-Up Dialog  
- **Functionality:**  
  - Allow entry of a filter text.
  - Toggle enable/disable status.
  - Confirm creation of the filter, which then:
    - Updates the filter list in the Filter Manager.
    - Triggers a regeneration of data (if in filtered mode).
    - Updates the filter panel display.
- **Access Points:**  
  - From the filter panel.
  - From the context menu in the text panel.
  - From the **Edit > Create Filter** option in the top menu bar.

#### 3.4.4 Top Menu Bar and Its Functions  
- **File Menu:**  
  - **Load File:** Opens file dialog to choose and load a text file.
  - **Reload File:** Re-read the file, reinitializing data caching and filtering.
  - **Save File:** Save a revised copy of the text file with the currently applied filters to a new file.
  - **Close File:** Close the open file and clear or save state as needed.
- **Edit Menu:**  
  - **Create Filter:** As outlined above.
  - **Disable All Filters:** Set all filters to `enabled = false` and update the display.
  - **Enable All Filters:** Set all filters to `enabled = true` and update the display.
  - **Export Filters:** Output the current filter settings (filter text, enabled state, order) to a file.
  - **Load Filters:** Read filter settings from an external file and update in-memory state and the Filter Panel.
- **Help Menu:**  
  - Provide user documentation, version information, and contact details.

---

## 4. Non-Functional Requirements

### 4.1 Performance  
- **Large File Support:** Optimize to efficiently process files with chunked loading and caching.
- **UI Responsiveness:** Ensure that the UI remains responsive during I/O and filtering operations using asynchronous task management.

### 4.2 Cross-Platform Compatibility  
- **Platforms:** Linux, macOS, and Windows.
- **Preferred Languages:**  
  - **C++ (C++17 or C++20)** or **Rust** (with appropriate Raylib bindings) to ensure high performance and seamless integration with Raylib.
- **Build Tools:**  
  - CMake (for C++) or Cargo (for Rust).

### 4.3 Maintainability and Extensibility  
- **Modular Architecture:**  
  - Clear separation of the data logic, filtering, state management, and the UI layer.
  - Use of dedicated managers/controllers for different concerns.
- **Ease of Future Updates:**  
  - New features (such as additional UI panels or new filtering behavior) can be integrated by updating the corresponding input, configuration, or UI modules.
- **Documentation:**  
  - Code and design should be thoroughly documented to facilitate future maintenance and development.

---

## 5. Architectural Diagram

Below is a high-level diagram emphasizing the updated UI components and interactions:

┌───────────────────────┐
               │     Top Menu Bar      │
               │ (File, Edit, Help)    │
               └─────────┬─────────────┘
                         │
           ┌─────────────┴─────────────┐
           │      Split Window         │
           │  ┌───────────────┐        │
           │  │  Text Panel   │        │
           │  │ (Top Panel)   │◄────┐  │
           │  └───────────────┘     │  │
           │         ▲              │  │
 Right-Click on Row     (Context  │  │
       opens Create     Menu)     │  │
           │                      │  │
           │         ┌────────────┴─────────┐
           │         │  Filter Panel        │
           │         │  (Bottom Panel)      │
           │         │ - Create, Delete,    │
           │         │   Update, Reorder    │
           │         └────────────┬─────────┘
           └──────────────────────┘
                         │
       (Event/Notification Manager)
                         │
                         ▼
              ┌─────────────────────┐
              │ State Manager /     │
              │ Controller          │
              │ - Data & Filters    │
              └─────────────────────┘

---

## 6. Proposed Technology Stack

- **Programming Language:**  
  - **Option 1:** C++ (modern C++17/C++20) with Raylib.  
  - **Option 2:** Rust (with raylib-rs bindings) for enhanced memory safety and concurrency.

- **GUI Library:**  
  - Raylib for rendering cross-platform graphics and handling window/menu interactions.

- **Build Tools:**  
  - CMake (for C++), Cargo (for Rust).

- **Optional Libraries:**  
  - Asynchronous libraries (e.g., Boost.Asio for C++ or Tokio for Rust) for background tasks.
  - Logging libraries (e.g., spdlog for C++ or log for Rust), configuration libraries, etc.

---

## 7. Milestones & Deliverables

1. **Design & Architecture Approval:**  
   - Finalize all architectural diagrams, including detailed UI and interaction flows.
   - Agree on technology stack and coding guidelines.

2. **Prototype Implementation:**  
   - Create Data Manager, basic Raylib UI with split panels, simple file load, and basic filtering.
   - Implement top menu bar with foundational File and Edit menu options.

3. **Full Feature Implementation:**  
   - Develop complete Filter Manager for creating, deleting, updating, and reordering filters.
   - Implement right-click context menu in the text panel to trigger filter creation pop-up.
   - Fully integrate top menu functions (Load, Reload, Save, Close; Create Filter, Disable/Enable All Filters, Export/Load Filters).
   - Complete implementation of managers/controllers for state, UI, event notifications, and input handling.
   - Add asynchronous I/O processing when required.

4. **Testing & Cross-Platform Validation:**  
   - Conduct tests on Linux, macOS, and Windows for performance, memory, and UI responsiveness.
   - Verify all menu options, panel interactions, and pop-up dialogs work as intended.

5. **Documentation & Deployment:**  
   - Provide comprehensive developer and user documentation.
   - Package deployment builds for all target platforms.

---

## 8. Updating the ARD for Future Requirements

When adding new features or updating the application, follow these steps to ensure a consistent update to the ARD:

1. **Identify the Affected Module(s):**  
   - Determine whether the change affects the Data Manager, Filter Manager, UI Manager, or another component.

2. **Describe the New Functionality:**  
   - Clearly define the new feature (e.g., a new panel, a new menu option, or additional filtering logic).

3. **Update the Functional Requirements Section:**  
   - Add or modify sections to include the new interactions and expected behavior.
   - Include details on interactions between the new functions and existing modules.

4. **Revise the Architectural Diagram:**  
   - Update diagrams to reflect any changes in the data flow or interactions among components.

5. **Document the User Interface Changes:**  
   - Update the GUI section with new menu items, panels, dialogues, or events.
   - Specify any additional event hooks or UI element behaviors.

6. **Define Acceptance Criteria and Test Cases:**  
   - Provide clear criteria for verifying that the updated functionality works as expected.

7. **Schedule Milestones and Deliverables:**  
   - Adjust project timelines to include integration and testing of the new features.

---

## 9. Summary

This updated ARD outlines a cross-platform text analysis application built with Raylib as the GUI framework and designed using a modular architecture. The application will:

- Load and process very large text files efficiently using chunked I/O and caching.
- Support dynamic filtering with configurable filters that can be created, edited, reordered, and managed through both a dedicated panel and contextual interactions.
- Provide a GUI with a split-window layout, top menu bar with File, Edit, and Help menus, and pop-up dialogues for filter management.
- Use a structured set of managers/controllers (Data Manager, Filter Manager, State Manager, UI Manager, Event Manager, etc.) to decouple core functionality from the presentation layer, ensuring both maintainability and extensibility.

This ARD serves as a blueprint for current development and as a guide for integrating new features or changes in the future.

---