#include <QtSingleApplication>
#include "slider.h"

int main(int argc, char *argv[]) {
    QtSingleApplication app(argc, argv);
    if(app.isRunning())
        return 0;
    Slider slider;
    return app.exec();
}
