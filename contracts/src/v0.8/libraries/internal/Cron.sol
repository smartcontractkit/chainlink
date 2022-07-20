// SPDX-License-Identifier: MIT

/*
  The Cron contract serves two primary functions:
    * parsing cron-formatted strings like "0 0 * * *" into
      structs called "Specs"
    * computing the "next tick" of a cron spec

  Because manipulating strings is gas-expensive in solidity,
  the intended use of this contract is for users to first convert
  their cron strings to encoded Spec structs via toEncodedSpec().
  Then, the user stores the Spec on chain. Finally, users use the nextTick(),
  function to determine the datetime of the next cron job run.

  Cron jobs are interpreted according to this format:

  ┌───────────── minute (0 - 59)
  │ ┌───────────── hour (0 - 23)
  │ │ ┌───────────── day of the month (1 - 31)
  │ │ │ ┌───────────── month (1 - 12)
  │ │ │ │ ┌───────────── day of the week (0 - 6) (Monday to Sunday)
  │ │ │ │ │
  │ │ │ │ │
  │ │ │ │ │
  * * * * *

  Special limitations:
    * there is no year field
    * no special characters: ? L W #
    * lists can have a max length of 26
    * no words like JAN / FEB or MON / TUES
*/

pragma solidity 0.8.6;

import "../../vendor/Strings.sol";
import "../../vendor/DateTime.sol";

// The fields of a cron spec, by name
string constant MINUTE = "minute";
string constant HOUR = "hour";
string constant DAY = "day";
string constant MONTH = "month";
string constant DAY_OF_WEEK = "day of week";

error UnknownFieldType();
error InvalidSpec(string reason);
error InvalidField(string field, string reason);
error ListTooLarge();

// Set of enums representing a cron field type
enum FieldType {
  WILD,
  EXACT,
  INTERVAL,
  RANGE,
  LIST
}

// A spec represents a cron job by decomposing it into 5 fields
struct Spec {
  Field minute;
  Field hour;
  Field day;
  Field month;
  Field dayOfWeek;
}

// A field represents a single element in a cron spec. There are 5 types
// of fields (see above). Not all properties of this struct are present at once.
struct Field {
  FieldType fieldType;
  uint8 singleValue;
  uint8 interval;
  uint8 rangeStart;
  uint8 rangeEnd;
  uint8 listLength;
  uint8[26] list;
}

/**
 * @title The Cron library
 * @notice A utility contract for encoding/decoding cron strings (ex: 0 0 * * *) into an
 * abstraction called a Spec. The library also includes a spec function, nextTick(), which
 * determines the next time a cron job should fire based on the current block timestamp.
 */
library Cron {
  using strings for *;

  /**
   * @notice nextTick calculates the next datetime that a spec "ticks", starting
   * from the current block timestamp. This is gas-intensive and therefore should
   * only be called off-chain.
   * @param spec the spec to evaluate
   * @return the next tick
   * @dev this is the internal version of the library. There is also an external version.
   */
  function nextTick(Spec memory spec) internal view returns (uint256) {
    uint16 year = DateTime.getYear(block.timestamp);
    uint8 month = DateTime.getMonth(block.timestamp);
    uint8 day = DateTime.getDay(block.timestamp);
    uint8 hour = DateTime.getHour(block.timestamp);
    uint8 minute = DateTime.getMinute(block.timestamp);
    uint8 dayOfWeek;
    for (; true; year++) {
      for (; month <= 12; month++) {
        if (!matches(spec.month, month)) {
          day = 1;
          hour = 0;
          minute = 0;
          continue;
        }
        uint8 maxDay = DateTime.getDaysInMonth(month, year);
        for (; day <= maxDay; day++) {
          if (!matches(spec.day, day)) {
            hour = 0;
            minute = 0;
            continue;
          }
          dayOfWeek = DateTime.getWeekday(DateTime.toTimestamp(year, month, day));
          if (!matches(spec.dayOfWeek, dayOfWeek)) {
            hour = 0;
            minute = 0;
            continue;
          }
          for (; hour < 24; hour++) {
            if (!matches(spec.hour, hour)) {
              minute = 0;
              continue;
            }
            for (; minute < 60; minute++) {
              if (!matches(spec.minute, minute)) {
                continue;
              }
              return DateTime.toTimestamp(year, month, day, hour, minute);
            }
            minute = 0;
          }
          hour = 0;
        }
        day = 1;
      }
      month = 1;
    }
  }

  /**
   * @notice lastTick calculates the previous datetime that a spec "ticks", starting
   * from the current block timestamp. This is gas-intensive and therefore should
   * only be called off-chain.
   * @param spec the spec to evaluate
   * @return the next tick
   */
  function lastTick(Spec memory spec) internal view returns (uint256) {
    uint16 year = DateTime.getYear(block.timestamp);
    uint8 month = DateTime.getMonth(block.timestamp);
    uint8 day = DateTime.getDay(block.timestamp);
    uint8 hour = DateTime.getHour(block.timestamp);
    uint8 minute = DateTime.getMinute(block.timestamp);
    uint8 dayOfWeek;
    bool resetDay;
    for (; true; year--) {
      for (; month > 0; month--) {
        if (!matches(spec.month, month)) {
          resetDay = true;
          hour = 23;
          minute = 59;
          continue;
        }
        if (resetDay) {
          day = DateTime.getDaysInMonth(month, year);
        }
        for (; day > 0; day--) {
          if (!matches(spec.day, day)) {
            hour = 23;
            minute = 59;
            continue;
          }
          dayOfWeek = DateTime.getWeekday(DateTime.toTimestamp(year, month, day));
          if (!matches(spec.dayOfWeek, dayOfWeek)) {
            hour = 23;
            minute = 59;
            continue;
          }
          for (; hour >= 0; hour--) {
            if (!matches(spec.hour, hour)) {
              minute = 59;
              if (hour == 0) {
                break;
              }
              continue;
            }
            for (; minute >= 0; minute--) {
              if (!matches(spec.minute, minute)) {
                if (minute == 0) {
                  break;
                }
                continue;
              }
              return DateTime.toTimestamp(year, month, day, hour, minute);
            }
            minute = 59;
            if (hour == 0) {
              break;
            }
          }
          hour = 23;
        }
        resetDay = true;
      }
      month = 12;
    }
  }

  /**
   * @notice matches evaluates whether or not a spec "ticks" at a given timestamp
   * @param spec the spec to evaluate
   * @param timestamp the timestamp to compare against
   * @return true / false if they match
   */
  function matches(Spec memory spec, uint256 timestamp) internal view returns (bool) {
    DateTime._DateTime memory dt = DateTime.parseTimestamp(timestamp);
    return
      matches(spec.month, dt.month) &&
      matches(spec.day, dt.day) &&
      matches(spec.hour, dt.hour) &&
      matches(spec.minute, dt.minute);
  }

  /**
   * @notice toSpec converts a cron string to a spec struct. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param cronString the cron string
   * @return the spec struct
   */
  function toSpec(string memory cronString) internal pure returns (Spec memory) {
    strings.slice memory space = strings.toSlice(" ");
    strings.slice memory cronSlice = strings.toSlice(cronString);
    if (cronSlice.count(space) != 4) {
      revert InvalidSpec("4 spaces required");
    }
    strings.slice memory minuteSlice = cronSlice.split(space);
    strings.slice memory hourSlice = cronSlice.split(space);
    strings.slice memory daySlice = cronSlice.split(space);
    strings.slice memory monthSlice = cronSlice.split(space);
    // DEV: dayOfWeekSlice = cronSlice
    // The cronSlice now contains the last section of the cron job,
    // which corresponds to the day of week
    if (
      minuteSlice.len() == 0 ||
      hourSlice.len() == 0 ||
      daySlice.len() == 0 ||
      monthSlice.len() == 0 ||
      cronSlice.len() == 0
    ) {
      revert InvalidSpec("some fields missing");
    }
    return
      validate(
        Spec({
          minute: sliceToField(minuteSlice),
          hour: sliceToField(hourSlice),
          day: sliceToField(daySlice),
          month: sliceToField(monthSlice),
          dayOfWeek: sliceToField(cronSlice)
        })
      );
  }

  /**
   * @notice toEncodedSpec converts a cron string to an abi-encoded spec. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param cronString the cron string
   * @return the abi-encoded spec
   */
  function toEncodedSpec(string memory cronString) internal pure returns (bytes memory) {
    return abi.encode(toSpec(cronString));
  }

  /**
   * @notice toCronString converts a cron spec to a human-readable cron string. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param spec the cron spec
   * @return the corresponding cron string
   */
  function toCronString(Spec memory spec) internal pure returns (string memory) {
    return
      string(
        bytes.concat(
          fieldToBstring(spec.minute),
          " ",
          fieldToBstring(spec.hour),
          " ",
          fieldToBstring(spec.day),
          " ",
          fieldToBstring(spec.month),
          " ",
          fieldToBstring(spec.dayOfWeek)
        )
      );
  }

  /**
   * @notice matches evaluates if a values matches a field.
   * ex: 3 matches *, 3 matches 0-5, 3 does not match 0,2,4
   * @param field the field struct to match against
   * @param value the value of a field
   * @return true / false if they match
   */
  function matches(Field memory field, uint8 value) private pure returns (bool) {
    if (field.fieldType == FieldType.WILD) {
      return true;
    } else if (field.fieldType == FieldType.INTERVAL) {
      return value % field.interval == 0;
    } else if (field.fieldType == FieldType.EXACT) {
      return value == field.singleValue;
    } else if (field.fieldType == FieldType.RANGE) {
      return value >= field.rangeStart && value <= field.rangeEnd;
    } else if (field.fieldType == FieldType.LIST) {
      for (uint256 idx = 0; idx < field.listLength; idx++) {
        if (value == field.list[idx]) {
          return true;
        }
      }
      return false;
    }
    revert UnknownFieldType();
  }

  // VALIDATIONS

  /**
   * @notice validate validates a spec, reverting if any errors are found
   * @param spec the spec to validate
   * @return the original spec
   */
  function validate(Spec memory spec) private pure returns (Spec memory) {
    validateField(spec.dayOfWeek, DAY_OF_WEEK, 0, 6);
    validateField(spec.month, MONTH, 1, 12);
    uint8 maxDay = maxDayForMonthField(spec.month);
    validateField(spec.day, DAY, 1, maxDay);
    validateField(spec.hour, HOUR, 0, 23);
    validateField(spec.minute, MINUTE, 0, 59);
    return spec;
  }

  /**
   * @notice validateField validates the value of a field. It reverts if an error is found.
   * @param field the field to validate
   * @param fieldName the name of the field ex "minute" or "hour"
   * @param min the minimum value a field can have (usually 1 or 0)
   * @param max the maximum value a field can have (ex minute = 59, hour = 23)
   */
  function validateField(
    Field memory field,
    string memory fieldName,
    uint8 min,
    uint8 max
  ) private pure {
    if (field.fieldType == FieldType.WILD) {
      return;
    } else if (field.fieldType == FieldType.EXACT) {
      if (field.singleValue < min || field.singleValue > max) {
        string memory reason = string(
          bytes.concat("value must be >=,", uintToBString(min), " and <=", uintToBString(max))
        );
        revert InvalidField(fieldName, reason);
      }
    } else if (field.fieldType == FieldType.INTERVAL) {
      if (field.interval < 1 || field.interval > max) {
        string memory reason = string(
          bytes.concat("inverval must be */(", uintToBString(1), "-", uintToBString(max), ")")
        );
        revert InvalidField(fieldName, reason);
      }
    } else if (field.fieldType == FieldType.RANGE) {
      if (field.rangeEnd > max || field.rangeEnd <= field.rangeStart) {
        string memory reason = string(
          bytes.concat("inverval must be within ", uintToBString(min), "-", uintToBString(max))
        );
        revert InvalidField(fieldName, reason);
      }
    } else if (field.fieldType == FieldType.LIST) {
      if (field.listLength < 2) {
        revert InvalidField(fieldName, "lists must have at least 2 items");
      }
      string memory reason = string(
        bytes.concat("items in list must be within ", uintToBString(min), "-", uintToBString(max))
      );
      uint8 listItem;
      for (uint256 idx = 0; idx < field.listLength; idx++) {
        listItem = field.list[idx];
        if (listItem < min || listItem > max) {
          revert InvalidField(fieldName, reason);
        }
      }
    } else {
      revert UnknownFieldType();
    }
  }

  /**
   * @notice maxDayForMonthField returns the maximum valid day given the month field
   * @param month the month field
   * @return the max day
   */
  function maxDayForMonthField(Field memory month) private pure returns (uint8) {
    // DEV: ranges are always safe because any two consecutive months will always
    // contain a month with 31 days
    if (month.fieldType == FieldType.WILD || month.fieldType == FieldType.RANGE) {
      return 31;
    } else if (month.fieldType == FieldType.EXACT) {
      // DEV: assume leap year in order to get max value
      return DateTime.getDaysInMonth(month.singleValue, 4);
    } else if (month.fieldType == FieldType.INTERVAL) {
      if (month.interval == 9 || month.interval == 11) {
        return 30;
      } else {
        return 31;
      }
    } else if (month.fieldType == FieldType.LIST) {
      uint8 result;
      for (uint256 idx = 0; idx < month.listLength; idx++) {
        // DEV: assume leap year in order to get max value
        uint8 daysInMonth = DateTime.getDaysInMonth(month.list[idx], 4);
        if (daysInMonth == 31) {
          return daysInMonth;
        }
        if (daysInMonth > result) {
          result = daysInMonth;
        }
      }
      return result;
    } else {
      revert UnknownFieldType();
    }
  }

  /**
   * @notice sliceToField converts a strings.slice to a field struct
   * @param fieldSlice the slice of a string representing the field of a cron job
   * @return the field
   */
  function sliceToField(strings.slice memory fieldSlice) private pure returns (Field memory) {
    strings.slice memory star = strings.toSlice("*");
    strings.slice memory dash = strings.toSlice("-");
    strings.slice memory slash = strings.toSlice("/");
    strings.slice memory comma = strings.toSlice(",");
    Field memory field;
    if (fieldSlice.equals(star)) {
      field.fieldType = FieldType.WILD;
    } else if (fieldSlice.contains(dash)) {
      field.fieldType = FieldType.RANGE;
      strings.slice memory start = fieldSlice.split(dash);
      field.rangeStart = sliceToUint8(start);
      field.rangeEnd = sliceToUint8(fieldSlice);
    } else if (fieldSlice.contains(slash)) {
      field.fieldType = FieldType.INTERVAL;
      fieldSlice.split(slash);
      field.interval = sliceToUint8(fieldSlice);
    } else if (fieldSlice.contains(comma)) {
      field.fieldType = FieldType.LIST;
      strings.slice memory token;
      while (fieldSlice.len() > 0) {
        if (field.listLength > 25) {
          revert ListTooLarge();
        }
        token = fieldSlice.split(comma);
        field.list[field.listLength] = sliceToUint8(token);
        field.listLength++;
      }
    } else {
      // needs input validation
      field.fieldType = FieldType.EXACT;
      field.singleValue = sliceToUint8(fieldSlice);
    }
    return field;
  }

  /**
   * @notice fieldToBstring converts a field to the bytes representation of that field string
   * @param field the field to stringify
   * @return bytes representing the string, ex: bytes("*")
   */
  function fieldToBstring(Field memory field) private pure returns (bytes memory) {
    if (field.fieldType == FieldType.WILD) {
      return "*";
    } else if (field.fieldType == FieldType.EXACT) {
      return uintToBString(uint256(field.singleValue));
    } else if (field.fieldType == FieldType.RANGE) {
      return bytes.concat(uintToBString(field.rangeStart), "-", uintToBString(field.rangeEnd));
    } else if (field.fieldType == FieldType.INTERVAL) {
      return bytes.concat("*/", uintToBString(uint256(field.interval)));
    } else if (field.fieldType == FieldType.LIST) {
      bytes memory result = uintToBString(field.list[0]);
      for (uint256 idx = 1; idx < field.listLength; idx++) {
        result = bytes.concat(result, ",", uintToBString(field.list[idx]));
      }
      return result;
    }
    revert UnknownFieldType();
  }

  /**
   * @notice uintToBString converts a uint256 to a bytes representation of that uint as a string
   * @param n the number to stringify
   * @return bytes representing the string, ex: bytes("1")
   */
  function uintToBString(uint256 n) private pure returns (bytes memory) {
    if (n == 0) {
      return "0";
    }
    uint256 j = n;
    uint256 len;
    while (j != 0) {
      len++;
      j /= 10;
    }
    bytes memory bstr = new bytes(len);
    uint256 k = len;
    while (n != 0) {
      k = k - 1;
      uint8 temp = (48 + uint8(n - (n / 10) * 10));
      bytes1 b1 = bytes1(temp);
      bstr[k] = b1;
      n /= 10;
    }
    return bstr;
  }

  /**
   * @notice sliceToUint8 converts a strings.slice to uint8
   * @param slice the string slice to convert to a uint8
   * @return the number that the string represents ex: "20" --> 20
   */
  function sliceToUint8(strings.slice memory slice) private pure returns (uint8) {
    bytes memory b = bytes(slice.toString());
    uint8 i;
    uint8 result = 0;
    for (i = 0; i < b.length; i++) {
      uint8 c = uint8(b[i]);
      if (c >= 48 && c <= 57) {
        result = result * 10 + (c - 48);
      }
    }
    return result;
  }
}
