#include "required/swift/required.djinni.swiftxx.h"

#include <iostream>
#include <string>

int main() {
  ::required::v1::Req message;
  message.set_count(7);

  try {
    (void)::djinni::swift::required::v1::Req::Translator::fromCpp(message);
    std::cerr << "fromCpp accepted a message without its required name field\n";
    return 1;
  } catch (const ::djinni::swift::ErrorValue& error) {
    if (std::string(error.what()).find("name") == std::string::npos) {
      std::cerr << "missing-field error did not name the required field: " << error.what() << '\n';
      return 1;
    }
  }

  message.set_name("ready");
  const auto value = ::djinni::swift::required::v1::Req::Translator::fromCpp(message);
  const auto roundTrip = ::djinni::swift::required::v1::Req::Translator::toCpp(value);
  if (!roundTrip.IsInitialized() || roundTrip.name() != "ready" || roundTrip.count() != 7) {
    std::cerr << "round trip did not preserve the initialized message\n";
    return 1;
  }

  std::cout << "required Swift translator rejects an uninitialized message and round trips an initialized message\n";
  return 0;
}
