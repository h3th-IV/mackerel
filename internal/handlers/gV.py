import random
import csv
import sys
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
import time
import os

#func --generate random ID numbers
def generate_random_ids(num_ids, id_length=11):
    generated_ids = []

    for _ in range(num_ids):
        while True:
            first_digit = random.choice('123456789')
            remaining_digits = ''.join(random.choices('0123456789', k=id_length - 1))
            new_id = first_digit + remaining_digits

            if validate_id(new_id):
                generated_ids.append(new_id)
                break

    return generated_ids

#func --generate serial ID numbers
def generte_serial_ids(start_id, num_ids):
    generated_ids = []
    current_id = int(start_id)

    for _ in range(num_ids):
        while True:
            new_id = str(current_id).zfill(len(start_id))  #ensure correct length with leading zeros
            if validate_id(new_id):
                generated_ids.append(new_id)
                break
        increment = random.randint(1000, 9999)  #incremet byrandom 4-digit nu
        current_id += increment

    return generated_ids

#validation func for gen IDs
def validate_id(id_number):
    max_zero_seq = 2
    max_non_zero_seq = 3
    zero_count = 0
    non_zero_count = 1

    prev_digit = id_number[0]

    for digit in id_number[1:]:
        if digit == '0':
            zero_count += 1
            non_zero_count = 0
        else:
            if digit == prev_digit:
                non_zero_count += 1
            else:
                non_zero_count = 1
            zero_count = 0

        if zero_count > max_zero_seq or non_zero_count > max_non_zero_seq:
            return False

        prev_digit = digit

    return True


#function to validate generated ID with api using selenium
def validate_on_website(id_number):
    current_dir = os.getcwd()
    service = Service(current_dir+'/chromedriver')
    driver = webdriver.Chrome(service=service)
    
    driver.get("https://smartkyc.ikejaelectric.com/nin_meter_link.xhtml")
    
    nin_field = driver.find_element(By.ID, "nin")
    nin_field.send_keys(id_number)
    
    verify_field = driver.find_element(By.ID, "nincode")
    verify_field.click()
    
    time.sleep(20)
    
    try:
        error_message = driver.find_element(By.CLASS_NAME, "ui-messages-error-summary")
        if "Record not found" in error_message.text:
            driver.quit()
            return "Invalid"
    except:
        pass

    driver.quit()
    return "Valid"


#save the generated ID numbers and their validation status to a CSV file
def save_to_csv(id_numbers, filename='validation_results.csv'):
    with open(filename, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(["ID Number", "Status"])
        for id_number, status in id_numbers:
            writer.writerow([id_number, status])

def main(generation_method, num_ids):
    if generation_method == 'serial':
        start_id = "30331470319"  #start with a known idd
        generated_ids = generte_serial_ids(start_id, num_ids)
    elif generation_method == 'random':
        generated_ids = generate_random_ids(num_ids)
    else:
        print("invalid generation method. choose 'serial' or 'random'.")
        return

    results = []
    #do validtion here
    for id_number in generated_ids:
        status = validate_on_website(id_number)
        results.append((id_number, status))
        print(f"{id_number}: {status}")
    
    save_to_csv(results)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python ninGenValidator.py <generation method: serial|random> <number_of_ids>")
    else:
        try:
            generation_method = sys.argv[1]
            num_ids = int(sys.argv[2])
            main(generation_method, num_ids)
        except ValueError:
            print("Please provide a valid integer for the number of IDs to generate.")
