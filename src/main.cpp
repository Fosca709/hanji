#include <QApplication>
#include <QIcon>
#include <QPlainTextEdit>

int main(int argc, char *argv[])
{
    QApplication application(argc, argv);
    application.setApplicationName(QStringLiteral("Hanji"));
    application.setWindowIcon(QIcon(QStringLiteral(":/assets/icon.svg")));

    QPlainTextEdit editor;
    editor.setWindowTitle(QStringLiteral("Hanji"));
    editor.setWindowIcon(QIcon(QStringLiteral(":/assets/icon.svg")));
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
    ));
    editor.show();

    return application.exec();
}
