import QtQuick
import QtQuick.Controls
import QtQuick.Window

Window {
    id: appWindow
    width: 350
    height: 350
    visible: true
    title: "Hanji"
    color: "#F8E58C"
    flags: alwaysOnTopMenuItem.checked ? Qt.Window | Qt.WindowStaysOnTopHint : Qt.Window

    Shortcut {
        sequence: "Ctrl+Q"
        onActivated: alwaysOnTopMenuItem.toggle()
    }

    ScrollView {
        anchors.fill: parent
        padding: 0

        TextArea {
            id: textArea
            width: parent.width
            wrapMode: TextArea.Wrap
            font.pointSize: 11
            color: "#2B1F0E"
            background: Rectangle {
                color: "#F8E58C"
            }

            TapHandler {
                acceptedButtons: Qt.RightButton

                onSingleTapped: (eventPoint, button) => {
                    textAreaContextMenu.x = eventPoint.position.x
                    textAreaContextMenu.y = eventPoint.position.y
                    textAreaContextMenu.open()
                }
            }

            Menu {
                id: textAreaContextMenu
                width: 130

                MenuItem {
                    id: alwaysOnTopMenuItem
                    text: "Always on Top"
                    checkable: true
                    leftPadding: 8
                    rightPadding: 8
                }

                MenuSeparator {}

                MenuItem {
                    text: "Undo"
                    enabled: textArea.canUndo
                    leftPadding: 8
                    rightPadding: 8
                    onTriggered: textArea.undo()
                }

                MenuItem {
                    text: "Redo"
                    enabled: textArea.canRedo
                    leftPadding: 8
                    rightPadding: 8
                    onTriggered: textArea.redo()
                }
            }
        }
    }
}
