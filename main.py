import sys
from pathlib import Path

from PySide6.QtGui import QGuiApplication, QIcon
from PySide6.QtQml import QQmlApplicationEngine


def main() -> int:
    app = QGuiApplication(sys.argv)
    icon_file = Path(__file__).resolve().parent / "assets" / "icon.svg"
    app.setWindowIcon(QIcon(str(icon_file)))
    engine = QQmlApplicationEngine()

    qml_file = Path(__file__).resolve().parent / "qml" / "Main.qml"
    engine.load(str(qml_file))

    if not engine.rootObjects():
        return 1

    return app.exec()


if __name__ == "__main__":
    raise SystemExit(main())
