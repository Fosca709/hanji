#include <QApplication>
#include <QAction>
#include <QContextMenuEvent>
#include <QDBusInterface>
#include <QIcon>
#include <QMenu>
#include <QPlainTextEdit>
#include <QTimer>

class NoteEditor final : public QPlainTextEdit
{
public:
    NoteEditor()
    {
        useKWinKeepAbove_ = QApplication::platformName().startsWith(QStringLiteral("wayland"))
            && qEnvironmentVariable("XDG_CURRENT_DESKTOP").contains(
                QStringLiteral("KDE"), Qt::CaseInsensitive);

        alwaysOnTopAction_.setText(QStringLiteral("Always on Top"));
        alwaysOnTopAction_.setShortcut(QKeySequence(QStringLiteral("Ctrl+Q")));
        alwaysOnTopAction_.setShortcutContext(Qt::WindowShortcut);
        updateAlwaysOnTopIcon();
        addAction(&alwaysOnTopAction_);

        connect(&alwaysOnTopAction_, &QAction::triggered, this,
                [this] { toggleAlwaysOnTop(!alwaysOnTop_); });
    }

protected:
    void contextMenuEvent(QContextMenuEvent *event) override
    {
        QMenu *menu = createStandardContextMenu();
        QAction *firstAction = menu->actions().value(0, nullptr);
        menu->insertAction(firstAction, &alwaysOnTopAction_);
        menu->insertSeparator(firstAction);
        menu->exec(event->globalPos());
        delete menu;
    }

private:
    void toggleAlwaysOnTop(bool checked)
    {
        if (useKWinKeepAbove_) {
            QDBusInterface kwinShortcuts(
                QStringLiteral("org.kde.kglobalaccel"),
                QStringLiteral("/component/kwin"),
                QStringLiteral("org.kde.kglobalaccel.Component"));
            kwinShortcuts.call(
                QStringLiteral("invokeShortcut"),
                QStringLiteral("Window Above Other Windows"));

            QTimer::singleShot(0, this, [this, checked] {
                alwaysOnTop_ = checked;
                updateAlwaysOnTopIcon();
            });
            return;
        }

        const QPoint previousPosition = pos();
        const bool restoreFocus = hasFocus();

        setWindowFlag(Qt::WindowStaysOnTopHint, checked);
        alwaysOnTop_ = checked;
        updateAlwaysOnTopIcon();
        show();
        move(previousPosition);

        if (restoreFocus) {
            QTimer::singleShot(0, this, [this] {
                activateWindow();
                setFocus(Qt::ShortcutFocusReason);
            });
        }
    }

    void updateAlwaysOnTopIcon()
    {
        alwaysOnTopAction_.setIcon(QIcon::fromTheme(
            alwaysOnTop_ ? QStringLiteral("window-unpin")
                         : QStringLiteral("window-pin")));
    }

    QAction alwaysOnTopAction_{this};
    bool alwaysOnTop_ = false;
    bool useKWinKeepAbove_ = false;
};

int main(int argc, char *argv[])
{
    QApplication application(argc, argv);
    application.setApplicationName(QStringLiteral("Hanji"));

    QIcon applicationIcon;
    applicationIcon.addFile(QStringLiteral(":/assets/icon-16.png"), QSize(16, 16));
    applicationIcon.addFile(QStringLiteral(":/assets/icon-32.png"), QSize(32, 32));
    application.setWindowIcon(applicationIcon);

    NoteEditor editor;
    editor.setWindowTitle(QStringLiteral("Hanji"));
    editor.setWindowIcon(applicationIcon);
    editor.resize(350, 350);
    editor.setLineWrapMode(QPlainTextEdit::WidgetWidth);
    editor.setWordWrapMode(QTextOption::WrapAtWordBoundaryOrAnywhere);
    editor.setStyleSheet(QStringLiteral(
        "QPlainTextEdit {"
        "  background-color: #F8E58C;"
        "  color: #000000;"
        "  border: 0;"
        "  padding: 6px;"
        "  font-size: 13pt;"
        "}"
        "QScrollBar:vertical {"
        "  background: transparent;"
        "  width: 12px;"
        "  margin: 4px 3px;"
        "}"
        "QScrollBar::handle:vertical {"
        "  background: rgba(112, 91, 39, 105);"
        "  border-radius: 3px;"
        "  min-height: 32px;"
        "}"
        "QScrollBar::handle:vertical:hover {"
        "  background: rgba(92, 72, 28, 150);"
        "}"
        "QScrollBar::handle:vertical:pressed {"
        "  background: rgba(72, 54, 18, 185);"
        "}"
        "QScrollBar::add-line:vertical,"
        "QScrollBar::sub-line:vertical {"
        "  height: 0;"
        "}"
        "QScrollBar::add-page:vertical,"
        "QScrollBar::sub-page:vertical {"
        "  background: transparent;"
        "}"
    ));
    editor.show();

    return application.exec();
}
